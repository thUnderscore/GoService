export TOOLCHAIN_DIR=armv7
export BIN_DIR=$ANDR_ARM7_DIR
export PKG_DIR=android_arm
export GOARCH=arm
export GOARM=7
export COMPILER=arm-linux-androideabi

sh build_android.sh
