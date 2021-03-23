#!/bin/bash

export ADDR=127.0.0.1
export PORT=60080
export PROT=tcp
export TIME=1

function usage() {
    echo "usage: smip.sh [-a addr] [-p port] [-P tcp|udp] [-t timeout]"
    echo "       smip -h"
}

while getopts "a:p:P:t:h" opt; do
    case $opt in
        a)  ADDR=$OPTARG ;;
        p)  PORT=$OPTARG ;;
        P)  PROT=$OPTARG ;;
        t)  TIME=$OPTARG ;;
        h)  usage; exit 1 ;;
        \?) usage; exit 1 ;;
    esac
done
shift $(($OPTIND - 1))

exec 5<>/dev/$PROT/$ADDR/$PORT
echo "" >&5
timeout $TIME cat <&5 2>/dev/null
