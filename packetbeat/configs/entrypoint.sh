#!/bin/bash
curl -XPUT -H 'Content-Type: application/json' http://localhost:9200/_template/packetbeat-7.0.0 -d@$PWD/configs/pd.json
./pbrunner