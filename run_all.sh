#!/bin/bash

DATE=$(date +"%m_%d_%y_%H%M")

echo $DATE

run_ql(){
    codeql query run -j -2 \
    -o bqrs/$1.bqrs \
    -d $2 \
    --additional-packs=/home/thevirg/.local/codeql-home/codeql-repo \
    --additional-packs=. \
    $3

codeql bqrs decode --format=json \
    -o json/$name.json \
    bqrs/$name.bqrs
}

# ## Run free5gc and get json
# name=free5gc_$DATE
# db=/home/thevirg/dev/ql_dbs/full_free5gc_with_deps
# query=go/service_access_discovery.ql
# run_ql $name $db $query

# ## Run SD-Core
# name=sdcore_$DATE
# db=/home/thevirg/dev/ql_dbs/sd_core_w_deps/go
# query=go/service_access_discovery.ql
# run_ql $name $db $query

# ## Run Open5gs Rel 16
# name=o5gs_rel16_$DATE
# db=/home/thevirg/dev/ql_dbs/o5gs_db
# query=open5gs/service_access_discovery.ql
# run_ql $name $db $query

# ## Run Open5gs Rel 17
# name=o5gs_rel17_$DATE
# db=/home/thevirg/dev/ql_dbs/o5gs_rel17
# query=open5gs/service_access_discovery.ql
# run_ql $name $db $query

## Run OAI
name=oai_$DATE
db=/home/thevirg/dev/ql_dbs/oai_db/cpp
query=OAI/service_access_discovery.ql
run_ql $name $db $query
