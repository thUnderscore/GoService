echo "build android/"$TOOLCHAIN_DIR

mkdir -p $BIN_DIR
cd $SRC_DIR

NDK=$NDK_TOOLCHAINS/$TOOLCHAIN_DIR
LIB_NAME=lib$PLUGIN_NAME.so

GOOS=android CGO_ENABLED=1 \
CC=$NDK/bin/$COMPILER-${MYCC} \
CXX=$NDK/bin/$COMPILER-${MYCCPP} \
CGO_CFLAGS="--sysroot=$NDK/sysroot" \
CGO_CPPFLAGS="--sysroot=$NDK/sysroot" \
CGO_LDFLAGS="--sysroot=$NDK/sysroot" \
$GOBIN/go build -o $LIB_NAME -buildmode=c-shared -pkgdir=$GOROOT/pkg/$PKG_DIR -tags=""

mv $LIB_NAME $BIN_DIR/$LIB_NAME

echo "build android/"$TOOLCHAIN_DIR" done: "$BIN_DIR/$LIB_NAME