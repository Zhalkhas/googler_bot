#!/bin/bash
echo "Starting reverse proxy"
/mercury-parser-api/reverse_proxy &
echo "Starting server"
cd /mercury-parser-api && yarn serve