#!/bin/bash

mkdir -p ~/.config/systemd/user

echo "
[Unit]
Description=HiddenLakeHelperComposite

[Service]
ExecStart=$HOME/.hidden-lake/bin/hlh-composite_amd64_linux --path $HOME/.hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden-lake-helper.composite.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hlh-composite_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hlh-composite_amd64_linux && \
    chmod +x hlh-composite_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden-lake-helper.composite.service
systemctl --user restart hidden-lake-helper.composite.service
