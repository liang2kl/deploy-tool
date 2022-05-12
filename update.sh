#!/bin/sh
(cd $1 && git checkout $3 && git pull)

if [ $? == "0" ]; then
    docker-compose -f $2/docker-compose.yml up -d --build $4
fi

exit $?
