#!/bin/bash

DATE=$(date +"%m_%d_%y_%H%M")
REPO_LOC=${REPO_LOC:-"."}
DB_LOC=${DB_LOC:-"$HOME"}
echo $DATE

run_ql(){

    echo $3

    codeql query run -j -2 \
    -o "${REPO_LOC}/bqrs/$1.bqrs" \
    -d "$2" \
    --additional-packs="$HOME/codeql-home/codeql-repo" \
    --additional-packs="${REPO_LOC}" \
    "$3"

    echo "Decoding query results to JSON."

codeql bqrs decode --format=json \
    -o ${REPO_LOC}/json/$1.json \
    ${REPO_LOC}/bqrs/$1.bqrs

    cd ${REPO_LOC}/gocode/Analyzer
    ./Analyzer -ap=${REPO_LOC}/json/$1.json

    mv "${REPO_LOC}/gocode/_output/acp.yaml" "${REPO_LOC}/output/$1.yaml"
    chmod 777 "${REPO_LOC}/output/$1.yaml"
    cd ${REPO_LOC}
}

cd ${REPO_LOC}/gocode/Analyzer
go build -o Analyzer main.go
cd ${REPO_LOC}

# ## Run free5gc and get json
name_f5gc=free5gc_$DATE
#db_f5gc=$HOME/dev/ql_dbs/full_free5gc_with_deps
db_f5gc=${DB_LOC}/free5gc/go
query_f5gc=go/service_access_discovery.ql
run_ql $name_f5gc $db_f5gc $query_f5gc

# ## Run SD-Core
sd_core_name=sd_core_$DATE
#sd_core_db=$HOME/dev/ql_dbs/sd_core_w_deps/go
sd_core_db=${DB_LOC}/sd_core/go
sd_core_query=go/service_access_discovery.ql
run_ql $sd_core_name $sd_core_db $sd_core_query

# ## Run Open5gs Rel 16
o5gs_r16_name=open5gs_r16_$DATE
#o5gs_r16_db=$HOME/dev/ql_dbs/o5gs_db
o5gs_r16_db=${DB_LOC}/open5gs_r16
o5gs_r16_query=open5gs/service_access_discovery.ql
run_ql $o5gs_r16_name $o5gs_r16_db $o5gs_r16_query

# ## Run Open5gs Rel 17
o5gs_r17_name=open5gs_r17_$DATE
#o5gs_r17_db=$HOME/dev/ql_dbs/o5gs_rel17
o5gs_r17_db=${DB_LOC}/open5gs_r17
o5gs_r17_query=open5gs/service_access_discovery.ql
run_ql $o5gs_r17_name $o5gs_r17_db $o5gs_r17_query

## Run OAI
oai_name=oai_$DATE
#oai_db=$HOME/dev/ql_dbs/oai_db/cpp
oai_db=${DB_LOC}/oai/cpp
oai_query=OAI/service_access_discovery.ql
run_ql $oai_name $oai_db $oai_query