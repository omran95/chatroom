package common

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var rateLimitRedisKeyPrefix = "chat:ratelimit"

var rateLimiterScript = `
local tokenBucketKey = KEYS[1]
local timestampKey = KEYS[2]

local fillingRate = tonumber(ARGV[1])
local bucketCapacity = tonumber(ARGV[2])

local currentTime = tonumber(ARGV[3])

local requestedTokens = tonumber(ARGV[4])

local expirationSeconds = math.floor(tonumber(ARGV[5]))

-- Check if timestamp key exists
local lastRefreshTime = tonumber(redis.call("get", timestampKey))
if lastRefreshTime == nil then
  lastRefreshTime = 0 -- Fallback for new users
end

local remainingTokens = tonumber(redis.call("get", tokenBucketKey))
if remainingTokens == nil then
    remainingTokens = bucketCapacity
end

local elapsedTime = math.max(0, currentTime - lastRefreshTime)
local refillableTokens = math.min(bucketCapacity, remainingTokens + (elapsedTime * fillingRate))

local allowedRequest = refillableTokens >= requestedTokens
local remainingTokensAfterRequest = refillableTokens - requestedTokens

if allowedRequest then
  redis.call("setex", tokenBucketKey, expirationSeconds, remainingTokensAfterRequest)
  redis.call("setex", timestampKey, expirationSeconds, currentTime)
end

return { allowedRequest }
`

type RateLimiter struct {
	redisClient    redis.UniversalClient
	FillingRate    float64
	bucketCapacity int
	expiration     time.Duration
	scriptSHA      string
}

func NewRateLimiter(redisClient redis.UniversalClient, fillingRate float64, bucketCapacity int, expiration time.Duration) (*RateLimiter, error) {
	scriptSHA, err := redisClient.ScriptLoad(context.Background(), rateLimiterScript).Result()
	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		redisClient:    redisClient,
		FillingRate:    fillingRate,
		bucketCapacity: bucketCapacity,
		expiration:     expiration,
		scriptSHA:      scriptSHA,
	}, nil

}

func (rateLimiter *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	tokensRequired := 1
	formattedKey := JoinStrings(rateLimitRedisKeyPrefix, ":", key)
	tokenBucketKey := JoinStrings("{", formattedKey, "}", ":tokens")
	timestampKey := JoinStrings("{", formattedKey, "}", ":ts")

	response, err := rateLimiter.redisClient.EvalSha(ctx, rateLimiter.scriptSHA, []string{tokenBucketKey, timestampKey}, rateLimiter.FillingRate, rateLimiter.bucketCapacity, time.Now().Unix(), tokensRequired, rateLimiter.expiration.Seconds()).Result()
	if err != nil {
		return false, err
	}

	result, _ := response.([]interface{})

	return result[0] == int64(1), nil
}

func JoinStrings(strs ...string) string {
	var stringBuilder strings.Builder
	for _, str := range strs {
		stringBuilder.WriteString(str)
	}
	return stringBuilder.String()
}
