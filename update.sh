#!/bin/sh
(cd $1 && git checkout $3 && git pull)

if [ $? == "0" ]; then
    docker-compose -f $2 up -d --force-recreate --build $4
fi

exit $?