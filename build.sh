#!/bin/sh

# You need to source the settings.sh file (source /path/to/UGE/default/common/settings.sh)
# of your Univa Grid Engine installation before building. "default" is the CELL name, which can
# be different in your setup.

if [ "$SGE_ROOT" = "" ]; then
    echo "source your Grid Engine settings.(c)sh file"
    exit 1
fi

ARCH=`$SGE_ROOT/util/arch`

export CGO_LDFLAGS="-L$SGE_ROOT/lib/$ARCH/"
export CGO_CFLAGS="-I$SGE_ROOT/include"
export LD_LIBRARY_PATH=$SGE_ROOT/lib/lx-amd64

go build -a
go install 
go test
