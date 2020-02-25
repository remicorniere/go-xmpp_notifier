#!/bin/bash
echo "${INPUT_SERVER_IP}   ${INPUT_SERVER_DOMAIN}">>/etc/hosts
go run main.go "${INPUT_SERVER_IP}"
