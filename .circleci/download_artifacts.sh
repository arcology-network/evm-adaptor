#!/bin/bash

export CIRCLE_TOKEN='43de5d9653835b09959f309954525b4b9b28ec8a'

curl https://circleci.com/api/v1.1/project/github/arcology-network/urlarbitrator-engine/latest/artifacts?circle-token=$CIRCLE_TOKEN \
   | grep -o 'https://[^"]*' \
   | sed -e "s/$/?circle-token=$CIRCLE_TOKEN/" \
   | wget -O /usr/local/lib/liburlarbitrator.so -v -i -