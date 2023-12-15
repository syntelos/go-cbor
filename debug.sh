#!/bin/bash

function initialize {
    rm -f debug.ini

    for lno in $(egrep -n 'BREAKPOINT' *.go | awk -F: '{print $1":"$2}')
    do
	echo "b ${lno}" >> debug.ini
    done
    echo "# debug.ini"
    cat -n debug.ini
}

if initialize 
then
    dlv --init debug.ini test
else
    1>&2 echo "$0 error reading BREAKPOINT from test.go"
    exit 1
fi
