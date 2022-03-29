#!/bin/bash

# Script for running vet command tests.
# Must be run from the project's root.

set -e

cd tests/vet/

echo -e "Running vet tests\n"
for dir in $(find . -maxdepth 4 -mindepth 4 -type d);
do
	echo "    $dir"
	cd $dir
	../../../../../../thdl vet > stdout || true
	diff --color stdout.golden stdout
	rm stdout
	cd ../../../..
done

echo -e "\nAll \e[1;32mPASSED\e[0m!"
