#!/bin/bash

if lno=$(egrep -n BREAKPOINT cbor_test.go | awk -F: '{print $1}') &&[ -n "${lno}" ]&&[ 1 -lt "${lno}" ]
then
    echo "b cbor_test.go:${lno}" > debug.ini

    dlv --init debug.ini test
else
    1>&2 echo "$0 error reading BREAKPOINT from test.go"
    exit 1
fi
