#!/bin/bash -e


# Put us in the same dir as the script.
cd $(dirname $0)

echo
echo "Running lib_openshift generate"
echo "------------------------------"
../generate.py


echo
echo "Running lib_utils Unit Tests"
echo "----------------------------"
cd unit

for test in *.py; do
    echo
    echo "--------------------------------------------------------------------------------"
    echo
    echo "Running $test..."
    ./$test
done


echo
echo "Running lib_utils Integration Tests"
echo "-----------------------------------"
cd ../integration

for test in *.yml; do
    echo
    echo "--------------------------------------------------------------------------------"
    echo
    echo "Running $test..."
    ./$test -vvv
done

# Clean up this damn file
# TODO: figure out why this is being written and clean it up.
rm kube-manager-test.yaml
