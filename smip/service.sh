#!/bin/bash

LOC_BIN_PATH=$(dirname $0)/smip
SYS_BIN_PATH=/usr/bin/smip
SERVICE_NAME=smip

function install_service() {
    systemctl stop ${SERVICE_NAME} 2>/dev/null

    install $LOC_BIN_PATH $SYS_BIN_PATH
    cat > /lib/systemd/system/${SERVICE_NAME}.service << EOF
[Unit]
Description=${SERVICE_NAME}

[Service]
ExecStart=${SYS_BIN_PATH}
Restart=always
User=root
Group=root
Environment=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
WorkingDirectory=/

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable ${SERVICE_NAME}
    systemctl start ${SERVICE_NAME}
}

function uninstall_service() {
    systemctl disable ${SERVICE_NAME} 2>/dev/null
    systemctl stop ${SERVICE_NAME} 2>/dev/null
    rm -f /lib/systemd/system/${SERVICE_NAME}.service
    rm -f $SYS_BIN_PATH
    systemctl daemon-reload
}

function run() {
    ${LOC_BIN_PATH} "$@"
}

if [ "$1" == "install" ] || [ "$1" == "i" ]; then
    install_service
elif [ "$1" == "uninstall" ] || [ "$1" == "u" ]; then
    uninstall_service
elif [ "$1" == "run" ] || [ "$1" == "r" ]; then
    shift
    run "$@"
else
    echo "usage: $0 <install|uninstall|run>"
fi

