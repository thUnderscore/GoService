#required externa variables
if [ -z $ANDROID_NDK_ROOT ]; then
	echo "ANDROID_NDK_ROOT must be defined"
	exit 1;
fi

if [ -z $GOROOT_BOOTSTRAP ]; then
	if [ -z $GOROOT ]; then
		echo "GOROOT_BOOTSTRAP or GOROOT must be defined"
		exit 1;
	fi
	export GOROOT_BOOTSTRAP=$GOROOT
fi

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OS="windows"
ARCH="amd64"
GO_BOOTSTRAP="1.7.1"
BatchExtension="bat"
MYCC=gcc

BUILD_DIR=$SCRIPT_DIR
export GOROOT=$BUILD_DIR/goroot
#mkdir -p $GOROOT
#mkdir -p $GOROOT/src
#mkdir -p $GOROOT/pkg
#mkdir -p $GOROOT/bin
export GOPATH=$BUILD_DIR/gopath
mkdir -p $GOPATH
export GOBIN=$GOROOT/bin

TOOLCHAINS_DIR=$BUILD_DIR/ndk-toolchain
GOSRC=$GOROOT/src
GOPKG=$GOROOT/pkg

echo "Download Go source"
#cd ${BUILD_DIR} && curl -s -L http://storage.googleapis.com/golang/go${GO_BOOTSTRAP}.src.tar.gz | tar xz && mv $BUILD_DIR/go $BUILD_DIR/goroot && cd $GOSRC

cd $GOSRC


NDK_TOOLCHAIN=$TOOLCHAINS_DIR/armv7
mkdir -p $NDK_TOOLCHAIN
echo "build  android/arm7 toolchain"
#$ANDROID_NDK_ROOT/build/tools/make-standalone-toolchain.sh --toolchain=arm-linux-androideabi-4.9 --platform=android-15 --install-dir=$NDK_TOOLCHAIN
echo "build  android/arm7 pkg"
cd $GOSRC
CC_FOR_TARGET=$NDK_TOOLCHAIN/bin/arm-linux-androideabi-${MYCC} GOOS=android GOARCH=arm GOARM=7 CGO_ENABLED=1 ./make.${BatchExtension} --no-clean || exit 1

if [ -d "$GOPKG/android_arm_shared" ]; then
  mkdir -p $GOPKG/android_arm
  mv $GOPKG/android_arm_shared/* $GOPKG/android_arm
  rm -rf $GOPKG/android_arm_shared
fi


NDK_TOOLCHAIN=$TOOLCHAINS_DIR/x86
mkdir -p $NDK_TOOLCHAIN
echo "build  android/x86  toolchain"
#$ANDROID_NDK_ROOT/build/tools/make-standalone-toolchain.sh --toolchain=x86-4.9 --platform=android-15 --install-dir=$NDK_TOOLCHAIN
echo "build  android/x86 pkg"
cd $GOSRC
CC_FOR_TARGET=$NDK_TOOLCHAIN/bin/i686-linux-android-${MYCC} GOOS=android GOARCH=386 CGO_ENABLED=1 ./make.${BatchExtension} --no-clean || exit 1

if [ -d "$GOPKG/android_386_shared" ]; then
  mkdir -p $GOPKG/android_386
  mv $GOPKG/android_386_shared/* $GOPKG/android_386
  rm -rf $GOPKG/android_386_shared
fi


cd $GOSRC
GOOS=windows GOARCH=386 CGO_ENABLED=1 ./make.${BatchExtension} --no-clean || exit 1

cd $GOSRC
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 ./make.${BatchExtension} --no-clean || exit 1




read -rsn1	

#echo "Remove build directory"

#rm -rf ${BUILD_DIR}