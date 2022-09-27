#!/bin/bash

# Script for running gen command tests.
# Must be run from the project's root.

set -e

cd tests/gen/

echo -e "Running gen tests\n"
for dir in $(find . -maxdepth 4 -mindepth 4 -type d);
do
	echo "    $dir"
	cd $dir
	../../../../../../hdl gen -to-stdout test.vhd > stdout || true
	diff --color stdout.golden.vhd stdout
	rm stdout
	cd ../../../..
done

echo -e "\nAll \e[1;32mPASSED\e[0m!"
