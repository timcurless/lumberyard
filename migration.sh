#!/bin/bash
set +e

sleep 25

# Create Lumberyard Keyspace
echo "CREATE KEYSPACE lumberyard WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};" | cqlsh db

# Create a custom Stack type
echo "use lumberyard; CREATE TYPE stack ( id text, assets map<text, text>, update_ts text, created_ts text);" | cqlsh db

# Create pipelines TABLE
#echo "use lumberyard; CREATE TABLE pipelines ( id UUID, name text, description text, PRIMARY KEY (id));" | cqlsh db

# Create stages TABLE
#echo "use lumberyard; CREATE TABLE stages ( id UUID, pipeline_id UUID, name text, description text, type text, version int, payload text, PRIMARY KEY(id));" | cqlsh db

# Create Projects TABLE
echo "use lumberyard; CREATE TABLE IF NOT EXISTS projects ( id text, name text, email text, update_ts text, created_ts text, stacks set<stack> PRIMARY KEY(id));" | cqlsh db

# Create Stacks TABLE
#echo "use lumberyard; CREATE TABLE IF NOT EXISTS stacks ( id varchar, project_id varchar, asset_ids set<varchar>, update_ts text, created_ts text PRIMARY KEY(id));" | cqlsh db
