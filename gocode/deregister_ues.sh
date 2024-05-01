#!/bin/bash

i=1

while [ $i -le 200 ]; do
    number=$(printf "%03d" "$i")
    ./nr-cli imsi-208930000000${number} -e "ps-release-all"
    i=$((i + 1))
done


i=1

while [ $i -le 200 ]; do
    number=$(printf "%03d" "$i")
    ./nr-cli imsi-208930000000${number} -e "deregister switch-off"
    i=$((i + 1))
done