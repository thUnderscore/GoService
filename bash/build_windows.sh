echo "build windows/"$GOARCH

mkdir -p $BIN_DIR
cd $SRC_DIR

LIB_NAME=$PLUGIN_NAME.dll

GOOS=windows CGO_ENABLED=1 \
$GOBIN/go build -o $PLUGIN_NAME.a -buildmode=c-archive -pkgdir=$GOROOT/pkg/$PKG_DIR

gcc -shared $M -L. -o $BIN_DIR/$LIB_NAME -Wl,--whole-archive $PLUGIN_NAME.a -Wl,--allow-multiple-definition,--no-whole-archive -static -lstdc++ -lwinmm -lntdll -lWs2_32 

rm $PLUGIN_NAME.a

#mv ../go/go_c_archive.a $BIN_DIR/api.a
#mv ../go/go_c_archive.h $BIN_DIR/api.h

echo "build windows/"$GOARCH" done: "$BIN_DIR/$LIB_NAME