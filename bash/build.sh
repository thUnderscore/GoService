#sh build.sh package_in/goroot_src/folder pluginname platform1|platform2|platform3 [build_to_dir]

export SRC_DIR=$GOPATH/src/$1
export PLUGIN_NAME=$2
PLATFORMS=$3
BUILD_DIR=$4

export MYCC=gcc
export MYCCPP=g++

if [ -z $BUILD_DIR ]; then
	BUILD_DIR=$TOOLS_ROOT_DIR/builds
fi

export ANDR_ARM7_DIR=$BUILD_DIR/Android/libs/armeabi-v7a
export ANDR_x86_DIR=$BUILD_DIR/Android/libs/x86
export DESKTOP_x86_DIR=$BUILD_DIR/x86
export DESKTOP_x64_DIR=$BUILD_DIR/x86_64


echo "starting build "$SRC_DIR" to "$BUILD_DIR

export IFS="|"
for platform in $PLATFORMS; do
  sh build_$platform.sh 
done

echo Done
#read -rsn1	