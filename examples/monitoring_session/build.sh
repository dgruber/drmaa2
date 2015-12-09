#!/bin/sh

# You need to source the settings.sh file (source /path/to/UGE/default/common/settings.sh)
# of your Univa Grid Engine installation before building. "default" is the CELL name, which can
# be different in your setup.

if [ "$SGE_ROOT" = "" ]; then
   echo "Please source your Univa Grid Engine environment before!"
   exit 1
fi

export CGO_LDFLAGS="-L$SGE_ROOT/lib/lx-amd64/"
export CGO_CFLAGS="-I$SGE_ROOT/include"

go build 

