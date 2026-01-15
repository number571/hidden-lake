#!/bin/bash

BASE64_BODY="$(\
    echo -n 'bash[@s]-c[@s]echo 'hello, world' >> file.txt && cat file.txt' | \
    base64 -w 0 \
)";
PUSH_FORMAT='{
    "method":"POST",
    "host":"hls-remoter",
    "path":"/exec",
    "head":{
        "Password": "DpxJFjAlrs4HOWga0wk14mZqQSBo9DxK"
    },
    "body":"'${BASE64_BODY}'"
}';

d="$(date +%s)";
curl -i -X POST "http://localhost:7572/api/network/request?friend=Bob" --data "${PUSH_FORMAT}";
echo && echo "Request took $(($(date +%s)-d)) seconds";
