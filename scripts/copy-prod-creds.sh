#!/bin/bash

if [ "$CONFIG" == "prod" ]; then
  git clone git@github.com:koding/credential.git
  cp credential/config/main.prod.coffee config/
fi