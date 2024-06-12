#!/bin/sh

case "$1" in
    "fresh-start")
        docker-compose up -d --scale chat-room=2 --build;;
    "start")
        docker-compose up -d --scale chat-room=2;;
    "stop")
        docker-compose stop;;
    "clean")
        docker-compose down;;
    *)
        echo "invalid command, only start, stop and clean are supported"
        exit 1;;
esac