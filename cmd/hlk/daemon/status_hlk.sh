#!/bin/bash

watch -c SYSTEMD_COLORS=1 systemctl --user status -o cat hidden_lake_kernel.service
