SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
TOOLS_ROOT_DIR=$SCRIPT_DIR/..
export GOROOT=$TOOLS_ROOT_DIR/goroot
export GOPATH=$TOOLS_ROOT_DIR/gopath
export GOBIN=$GOROOT/bin
export SRC_DIR=$GOPATH/src/$1
cd $SRC_DIR

$GOBIN/go tool compile 

echo Done
read -rsn1	