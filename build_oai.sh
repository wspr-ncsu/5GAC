#!/bin/bash

omec_nfs=("amf" "ausf" "nef" "nrf" "nssf" "pcf" "nef" "smf" "udm" "udr")
for ((i=0; i<${#omec_nfs[@]}; i++)); 
do 
	cd /go/oai-cn5g-${omec_nfs[i]};
	./build/scripts/build_${omec_nfs[i]} -f -j -V;
	cd ..;
        #codeql database --add-repo /dbs/oai --language=c-cpp --command "build/scripts/build_${omec_nfs[i]} -f -j -V";
done
