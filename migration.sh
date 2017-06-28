#!/bin/bash
set +e

sleep 25

# Create Lumberyard Keyspace
echo "CREATE KEYSPACE lumberyard WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};" | cqlsh db

# Create pipelines TABLE
echo "use lumberyard; CREATE TABLE pipelines ( id UUID, name text, description text, PRIMARY KEY (id));" | cqlsh db
