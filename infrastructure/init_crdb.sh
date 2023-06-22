#! /bin/bash

pwd
ls -al

mkdir -p ./certs/ca

cockroach cert create-ca --certs-dir ./certs --ca-key ./certs/ca/ca.key
cockroach cert create-node localhost node_1 --certs-dir ./certs --ca-key ./certs/ca/ca.key
cockroach cert create-client root --certs-dir ./certs --ca-key ./certs/ca/ca.key
cockroach cert create-client dbuser --certs-dir ./certs --ca-key ./certs/ca/ca.key

cockroach start-single-node --certs-dir ./certs --advertise-addr localhost

sleep 3

cockroach sql --user root --certs-dir ./certs --file ./bootstrap_crdb.sql