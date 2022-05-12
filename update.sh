#!/bin/sh
(cd $1/$2 && git checkout $3 && git pull)

if [ $? == "0" ]; then
    docker-compose -f $1/deploy/docker-compose.yml up -d --build $4
fi

exit $?
