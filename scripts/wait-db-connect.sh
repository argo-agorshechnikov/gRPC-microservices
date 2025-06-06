#!/bin/sh

HOST=$1
PORT=$2

echo "Waiting for db at $HOST:$PORT..."

shift 2

while ! nc -z $HOST $PORT; do
    echo "DB is unavailable - sleeping"
    sleep 2
done

echo "DB is up - executing command"
exec "${@}"