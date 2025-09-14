#!/bin/bash

set -e

pm2 start /app/web/server.js --name clio-web --cwd /app/web -i ${PM2_INSTANCES} --no-daemon