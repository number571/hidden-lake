#!/bin/bash

journalctl --user -n 100000 -o cat -u 'hidden-lake-helper.composite.service' | \
    grep -E ".*service=HLK.*" | \
    grep -xv -E "(.*method=.*)|(.*type=(BRDCS|UNDEC|ENQRQ).*)"
