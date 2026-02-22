#!/bin/bash

echo "
[Unit]
Description=HiddenLakeServiceMessenger

[Service]
ExecStart=$HOME/.hidden-lake/bin/hls-messenger_amd64_linux --path $HOME/hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden_lake_service_messenger.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hls-messenger_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hls-messenger_amd64_linux && \
    chmod +x hls-messenger_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden_lake_service_messenger.service
systemctl --user restart hidden_lake_service_messenger.service
