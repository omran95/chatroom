CREATE KEYSPACE chatRoom WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
USE chatRoom;
CREATE TABLE rooms (
    id varint,
    PRIMARY KEY((id))
);
CREATE TABLE messages (
    id varint,
    event int,
    room_id varint,
    payload text,
    seen boolean,
    timestamp timestamp,
    PRIMARY KEY((room_id), id)
) WITH CLUSTERING ORDER BY (id DESC);