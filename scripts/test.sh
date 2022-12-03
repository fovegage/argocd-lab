#!/bin/bash

dialName="dial-proxy"
if [ -f "$dialName" ]; then
  rm -rf "$dialName"
  kiil -9 898
  # shellcheck disable=SC2105
  continue
fi

ls , echo "\n\n"
