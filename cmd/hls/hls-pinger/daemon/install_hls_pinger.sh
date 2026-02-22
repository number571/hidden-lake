#!/bin/bash

echo "
[Unit]
Description=HiddenLakeServicePinger

[Service]
ExecStart=$HOME/.hidden-lake/bin/hls-pinger_amd64_linux --path $HOME/hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden-lake-service.pinger.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hls-pinger_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hls-pinger_amd64_linux && \
    chmod +x hls-pinger_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden-lake-service.pinger.service
systemctl --user restart hidden-lake-service.pinger.service
