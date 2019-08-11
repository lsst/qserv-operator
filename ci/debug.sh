#!/bin/sh

set -e
set -x

TOKEN='-RtOuVrv58G58tKu9TAWwA'
#TOKEN='nSetU36eTO4nrWFnfz1qNQ'
JOBID='570046978'
curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "Accept: application/json" \
  -H "Travis-API-Version: 3" \
  -H "Authorization: token ${TOKEN}" \
  -d "{\"quiet\": true}" \
  https://api.travis-ci.org/job/${JOBID}/debug