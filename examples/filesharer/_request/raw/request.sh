#!/bin/bash

PUSH_FORMAT='{
    "receiver":"Bob",
    "req_data":{
        "method":"GET",
        "host":"hls-filesharer",
        "path":"/list?page=0"
    }
}';

curl -i -X POST http://localhost:8572/api/network/request --data "${PUSH_FORMAT}";
echo 
