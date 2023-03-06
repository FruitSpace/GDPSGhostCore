#!/bin/bash
CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath 
R=ghostcore
sudo docker build -t cr.yandex/crpr24jcqm2dno6qlm3b/$R . && docker push cr.yandex/crpr24jcqm2dno6qlm3b/$R

