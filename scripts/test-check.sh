#!/bin/bash

# Script for running check command tests.
# Must be run from the project's root.

set -e

cd tests/check/

echo -e "Running check tests\n"
for dir in $(find . -maxdepth 4 -mindepth 4 -type d);
do
	echo "    $dir"
	cd $dir
	../../../../../../thdl check > stdout || true
	diff --color stdout.golden stdout
	rm stdout
	cd ../../../..
done

echo -e "\nAll \e[1;32mPASSED\e[0m!"
