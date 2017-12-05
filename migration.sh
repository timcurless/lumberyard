#!/bin/bash
set +e

sleep 25

# Create Lumberyard Keyspace
echo "CREATE KEYSPACE lumberyard WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};" | cqlsh db

# Create pipelines TABLE
#echo "use lumberyard; CREATE TABLE pipelines ( id UUID, name text, description text, PRIMARY KEY (id));" | cqlsh db

# Create stages TABLE
#echo "use lumberyard; CREATE TABLE stages ( id UUID, pipeline_id UUID, name text, description text, type text, version int, payload text, PRIMARY KEY(id));" | cqlsh db

# Create Projects TABLE
echo "use lumberyard; CREATE TABLE projects ( id text, name text, email text, update_ts text, created_ts text, PRIMARY KEY(id));" | cqlsh db
