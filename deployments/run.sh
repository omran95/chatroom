#!/bin/sh

ROOM_SCALE=3
SUBSCRIBER_SCALE=3

# Construct scale options
SCALE_OPTIONS="--scale chat-room=$ROOM_SCALE --scale subscriber=$SUBSCRIBER_SCALE"

case "$1" in
    "fresh-start")
        docker-compose up -d --build $SCALE_OPTIONS;;
    "start")
        docker-compose up -d $SCALE_OPTIONS;;
    "stop")
        docker-compose stop
        exit 0;;
    "clean")
        docker-compose down
        exit 0;;
    *)
        echo "invalid command, only fresh-start, start, stop, and clean are supported"
        exit 1;;
esac