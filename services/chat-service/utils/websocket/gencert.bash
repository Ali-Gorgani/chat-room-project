#!/bin/bash

echo "creating server.key"
openssl genrsa -out services/chat-service/server.key 2048
openssl ecparam -genkey -name secp384r1 -out services/chat-service/server.key
echo "creating server.crt"
openssl req -new -x509 -sha256 -key services/chat-service/server.key -out services/chat-service/server.crt -batch -days 365