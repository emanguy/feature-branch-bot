#!/bin/bash

if [[ "$BOT_CONFIG_JSON" != "" ]]
then
  echo "$BOT_CONFIG_JSON" > ./bot-config.json
fi

java -jar ./FeatureBranchBot.jar