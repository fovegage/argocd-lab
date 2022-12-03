#!/bin/bash

dialName="dial-proxy"
if [ -f "$dialName" ]; then
  rm -rf "$dialName"
fi