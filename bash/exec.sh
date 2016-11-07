#sh build.sh package_in/goroot_src/folder pluginname platform1|platform2|platform3 [build_to_dir]

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export TOOLS_ROOT_DIR=$SCRIPT_DIR/..
export GOROOT=$TOOLS_ROOT_DIR/goroot
export GOPATH=$TOOLS_ROOT_DIR/gopath
export NDK_TOOLCHAINS=$TOOLS_ROOT_DIR/ndk-toolchain
export GOBIN=$GOROOT/bin

Cmd=$1
shift

if [[ $Cmd == *.sh ]] ;
then
    sh $Cmd $@
else
    $Cmd $@	
fi

