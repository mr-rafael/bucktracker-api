#!/bin/sh
if [ ! -f /app/.env ]; then
    cp /app/.env.example /app/.env
fi
exec /app/server "$@"
