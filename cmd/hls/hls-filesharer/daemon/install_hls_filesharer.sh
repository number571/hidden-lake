#!/bin/bash

echo "
[Unit]
Description=HiddenLakeServiceFilesharer

[Service]
ExecStart=$HOME/.hidden-lake/bin/hls-filesharer_amd64_linux --path $HOME/hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden-lake-service.filesharer.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hls-filesharer_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hls-filesharer_amd64_linux && \
    chmod +x hls-filesharer_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden-lake-service.filesharer.service
systemctl --user restart hidden-lake-service.filesharer.service
