#!/bin/bash

PUSH_FORMAT='{
    "receiver":"Bob",
    "req_data":{
        "method":"GET",
		"host":"hidden-lake-pinger",
		"path":"/ping"
    }
}';

d="$(date +%s)";
curl -i -X POST http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo && echo "Request took $(($(date +%s)-d)) seconds";
