#!/usr/bin/env bash
cd linux
upx sklabel_cutting_webservice_linux
cd ..
docker rmi -f petrjahoda/sklabel_cutting_webservice:latest
docker build -t petrjahoda/sklabel_cutting_webservice:latest .
docker push petrjahoda/sklabel_cutting_webservice:latest

docker rmi -f petrjahoda/sklabel_cutting_webservice:2020.3.3
docker build -t petrjahoda/sklabel_cutting_webservice:2020.3.3 .
docker push petrjahoda/sklabel_cutting_webservice:2020.3.3
