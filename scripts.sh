#!/usr/bin/env bash
docker run --name mysql -e MYSQL_ROOT_PASSWORD=root -p 3308:3306 -d mysql
