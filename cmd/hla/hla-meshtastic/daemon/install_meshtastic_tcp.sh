#!/bin/bash

mkdir -p ~/.config/systemd/user

echo "
[Unit]
Description=HiddenLakeAdapterMeshtastic

[Service]
ExecStart=$HOME/.hidden-lake/bin/hla-meshtastic_amd64_linux --path $HOME/hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden-lake-adapter.meshtastic.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hla-meshtastic_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hla-meshtastic_amd64_linux && \
    chmod +x hla-meshtastic_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden-lake-adapter.meshtastic.service
systemctl --user restart hidden-lake-adapter.meshtastic.service
