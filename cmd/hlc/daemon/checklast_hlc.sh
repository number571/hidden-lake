#!/bin/bash

journalctl -n 100000 -o cat -u 'hidden_lake_composite.service' | \
    grep -E ".*service=HLS.*" | \
    grep -xv -E "(.*method=.*)|(.*type=(BRDCS|UNDEC|ENQRQ).*)"
