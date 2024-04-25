#!/bin/sh

if [ "$#" -eq 0 ] || [ "$#" -gt 1 ]; then
    echo "Illegal number of arguments. Usage: mongo.sh [ install | configure | uninstall | start | stop ]"
    exit 0
fi

if test $1 = "install"; then   
    docker run -d -p 27017:27017 --name mongodb mongo:7 mongod --replSet rs0

    docker exec -it mongodb mongosh --eval "use admin; db.adminCommand({shutdown: 1,comment: 'Convert to cluster'})"
    docker exec -it mongodb mongosh --eval "rs.initiate(); rs.conf(); rs.status()"

    echo "MongoDB replicaset installed and running successfully!"
    exit 0
fi

if test $1 = "configure"; then   
    docker exec -it mongodb mongosh --eval "use fadb; db.createCollection('accounts'); db.createCollection('budgets')"
    echo "MongoDB configured successfully!"
    exit 0
fi

if test $1 = "uninstall"; then
    docker rm mongodb
    exit 0
fi

if test $1 = "start"; then
    docker start mongodb
    exit 0
fi

if test $1 = "stop"; then
    docker stop mongodb
    exit 0
fi


