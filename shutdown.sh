#!/bin/bash

docker-compose -f docker-compose-bloc-web.yml down -v

if [[ "$OSTYPE" == "darwin"* ]]; then
	# Mac OSX
    docker-compose -f docker-compose-bloc-server-mac.yml down -v
elif [[ "$OSTYPE" == "linux"* ]]; then
	# Linux
    docker-compose -f docker-compose-bloc-server-linux.yml down -v
fi
docker-compose down -v