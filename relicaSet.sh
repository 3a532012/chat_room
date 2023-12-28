#!/bin/bash

echo "Running setup script..."

mongod --replSet rs0 --bind_ip_all &

sleep 2

mongosh <<EOF
rs.initiate()
rs.add('db2')
rs.add('db3')

db.adminCommand({
   "setDefaultRWConcern" : 1,
   "defaultWriteConcern" : {
     "w" : "majority"
   },
   "defaultReadConcern" : { "level" : "majority" }
 })
rs.addArb('db4')
db.getMongo().setReadPref('secondaryPreferred')
EOF

echo "Setup script completed."

tail -f /dev/null
