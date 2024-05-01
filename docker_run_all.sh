#!/bin/sh

export PATH="${PATH}:/go/codeql"
export REPO_LOC="/repo"
export DB_LOC="/dbs"
echo $PATH
${REPO_LOC}/run_all.sh
