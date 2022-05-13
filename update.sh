#!/bin/sh
(cd "$1" && git checkout "$3" && git pull)

ret=$?

if [ $ret != "0" ]; then
    exit $ret
fi

docker-compose -f "$2" up -d --force-recreate --build "$4"

exit $?