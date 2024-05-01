#!/bin/bash

script_dir=""$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )""
free5gc_path="$script_dir/free5gc"
dir="holder"
one='$1'
two='$2'

# Get deps of all NFs
cd $free5gc_path/NFs/
for d in */; do
	command cd $d > /dev/null
	echo "Copying in Deps for $d"
	# remove old folders if we are doing it again
	if [ -d "github.com/" ] || [ -d "golang.org/" ] || [ -d "go.mongodb.org/" ] || [ -d "google.golang.org/" ] ||
	 [ -d "gopkg.in/" ] || [ -d "git.cs.nctu.edu.tw/" ];
	 then
		echo "Removing old dep folders..."
    rm -rf github.com/ golang.org/ go.mongodb.org/ google.golang.org/ gopkg.in/ git.cs.nctu.edu.tw/
 	fi
	command go mod vendor
	command cp -rT vendor/ .
	command rm -rf vendor/
	# remove go.mod and go.sum from internals
	for subdir in */; do
		command cd $subdir > /dev/null
		echo "Removing go.mod and go.sum from $subdir"
		command find . -name go.mod -type f -delete
		command find . -name go.sum -type f -delete
		command cd .. > /dev/null
	done
	command cd .. > /dev/null
done

# Return to free5gc root to do perl replace
command cd $free5gc_path/NFs

# Replace all imports except github.com/golang and github.com/free5gc/<current nf>
for d in */; do
	echo "Replacing Imports in $d"
	cd $d
	dir=${d%/}
	if [ $dir == "upf" ] ; then
		dir="go-upf"
	fi
#	find_pregex="([[:blank:]]|import[[:blank:]])\"(github.com|google.golang.org|git.cs.nctu.edu.tw)\/(?!free5gc\/$dir)(?!golang)"
	find_pregex="([[:blank:]]|import[[:blank:]])\"(github.com|google.golang.org|git.cs.nctu.edu.tw)\/(?!free5gc\/$dir)"
	perl_replace="$one\"github.com\/free5gc\/$dir\/$two\/"
	full_replace="'s/$find_pregex/$perl_replace/g'"
	for f in $(find . -name '*.go'); do
		eval perl -pi -e "$full_replace" $f
		#echo $f
	done
	echo "Updating go.mod"
	command sed -i 's/go 1.17/go 1.19/g' go.mod
	command go mod tidy
	echo "Leaving $d"
	cd ..
done


cd $free5gc_path
echo "Making new database"
codeql database create /dbs/free5gc --language=go --overwrite --db-cluster --command="make"


