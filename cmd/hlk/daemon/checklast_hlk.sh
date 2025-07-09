#!/bin/bash

journalctl -n 100000 -o cat -u 'hidden_lake_kernel.service' | \
    grep -xv -E "(.*method=.*)|(.*type=(BRDCS|UNDEC|ENQRQ).*)"
