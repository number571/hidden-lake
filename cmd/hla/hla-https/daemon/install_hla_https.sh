#!/bin/bash

echo "
[Unit]
Description=HiddenLakeAdapterHTTPS

[Service]
ExecStart=$HOME/.hidden-lake/bin/hla-https_amd64_linux --path $HOME/hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden-lake-adapter.https.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hla-https_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hla-https_amd64_linux && \
    chmod +x hla-https_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden-lake-adapter.https.service
systemctl --user restart hidden-lake-adapter.https.service
