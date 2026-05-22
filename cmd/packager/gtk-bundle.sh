#!/usr/bin/env bash
# GTK 依赖打包脚本
# 功能：自动发现并打包 GTK 应用所需的运行时依赖
# 用法：gtk-bundle.sh <appdir> <gtk-version> [webkit-version]

set -e

APPDIR="$1"
GTK_VERSION="$2"
WEBKIT_VERSION="$3"  # 可选: "4.0", "4.1", ""

if [ -z "$APPDIR" ] || [ -z "$GTK_VERSION" ]; then
    echo "Usage: $0 <appdir> <gtk-version> [webkit-version]"
    echo "  appdir: AppDir 目录路径"
    echo "  gtk-version: GTK 版本 (2, 3)"
    echo "  webkit-version: WebKit 版本 (4.0, 4.1, 可选)"
    exit 1
fi

APPDIR="$(realpath "$APPDIR")"

# 查找 pkg-config
if command -v pkgconf > /dev/null; then
    PKG_CONFIG="pkgconf"
elif command -v pkg-config > /dev/null; then
    PKG_CONFIG="pkg-config"
else
    echo "ERROR: pkg-config not found"
    exit 1
fi

# 查找库目录路径
find_lib_path() {
    local paths=(
        "$(dpkg-architecture -qDEB_HOST_MULTIARCH 2>/dev/null && echo "/usr/lib/$(dpkg-architecture -qDEB_HOST_MULTIARCH)")"
        "/usr/lib64"
        "/usr/lib"
    )
    for p in "${paths[@]}"; do
        [ -d "$p" ] && echo "$p" && return 0
    done
}

LIB_PATH="$(find_lib_path)"

# 获取 pkg-config 变量
pkg_var() {
    local var="$1" lib="$2" default="$3"
    local val="$("$PKG_CONFIG" --variable="$var" "$lib" 2>/dev/null)"
    echo "${val:-$default}"
}

# 复制库目录下的文件
copy_libs() {
    local pattern="$1" src_dir="$2"
    while IFS= read -r -d '' file; do
        local rel="${file#$LIB_PATH/}"
        local dst_dir="$APPDIR/usr/lib/$(dirname "$rel")"
        mkdir -p "$dst_dir"
        cp -a "$file" "$dst_dir/"
    done < <(find "$src_dir" \( -type l -o -type f \) -name "$pattern" -print0 2>/dev/null)
}

# 创建 AppRun hook
HOOK_DIR="$APPDIR/apprun-hooks"
HOOK_FILE="$HOOK_DIR/gtk-bundle.sh"
mkdir -p "$HOOK_DIR"

cat > "$HOOK_FILE" << 'HOOK_EOF'
#!/usr/bin/env bash
export APPDIR="${APPDIR:-"$(dirname "$(realpath "$0")")"}"
export GTK_DATA_PREFIX="$APPDIR"
export GDK_BACKEND=x11
export XDG_DATA_DIRS="$APPDIR/usr/share:/usr/share:$XDG_DATA_DIRS"
HOOK_EOF

# 复制 GLib schemas
echo "Copying GLib schemas..."
SCHEMAS_DIR="$(pkg_var schemasdir gio-2.0 /usr/share/glib-2.0/schemas)"
if [ -d "$SCHEMAS_DIR" ]; then
    mkdir -p "$APPDIR/$SCHEMAS_DIR"
    cp -a "$SCHEMAS_DIR"/*.xml "$APPDIR/$SCHEMAS_DIR/" 2>/dev/null || true
    glib-compile-schemas "$APPDIR/$SCHEMAS_DIR" 2>/dev/null || true
    echo "export GSETTINGS_SCHEMA_DIR=\"\$APPDIR/$SCHEMAS_DIR\"" >> "$HOOK_FILE"
fi

# 复制 GIRepository typelibs
echo "Copying GIRepository typelibs..."
TYPELIB_DIR="$(pkg_var typelibdir gobject-introspection-1.0 "$LIB_PATH/girepository-1.0")"
if [ -d "$TYPELIB_DIR" ]; then
    mkdir -p "$APPDIR/usr/lib/girepository-1.0"
    cp -a "$TYPELIB_DIR"/*.typelib "$APPDIR/usr/lib/girepository-1.0/" 2>/dev/null || true
    echo "export GI_TYPELIB_PATH=\"\$APPDIR/usr/lib/girepository-1.0\"" >> "$HOOK_FILE"
fi

# GTK 版本特定处理
case "$GTK_VERSION" in
    2)
        echo "Processing GTK 2 modules..."
        GTK_LIBDIR="$(pkg_var libdir gtk+-2.0 "$LIB_PATH")/gtk-2.0"
        GTK_BINARY_VERSION="$(pkg_var gtk_binary_version gtk+-2.0 "2.10.0")"

        # 复制 immodules
        IMMODULES_DIR="$GTK_LIBDIR/$GTK_BINARY_VERSION/immodules"
        if [ -d "$IMMODULES_DIR" ]; then
            mkdir -p "$APPDIR/usr/lib/gtk-2.0/$GTK_BINARY_VERSION/immodules"
            cp -a "$IMMODULES_DIR"/*.so "$APPDIR/usr/lib/gtk-2.0/$GTK_BINARY_VERSION/immodules/" 2>/dev/null || true
            if command -v gtk-query-immodules-2.0 > /dev/null; then
                gtk-query-immodules-2.0 > "$APPDIR/usr/lib/gtk-2.0/$GTK_BINARY_VERSION/immodules.cache" 2>/dev/null || true
            fi
        fi

        # 复制 printbackends
        PRINTBACKENDS_DIR="$GTK_LIBDIR/$GTK_BINARY_VERSION/printbackends"
        if [ -d "$PRINTBACKENDS_DIR" ]; then
            mkdir -p "$APPDIR/usr/lib/gtk-2.0/$GTK_BINARY_VERSION/printbackends"
            cp -a "$PRINTBACKENDS_DIR"/*.so "$APPDIR/usr/lib/gtk-2.0/$GTK_BINARY_VERSION/printbackends/" 2>/dev/null || true
        fi

        GTK_EXEC_PREFIX="$(pkg_var exec_prefix gtk+-2.0 /usr)"
        cat >> "$HOOK_FILE" << EOF
export GTK_EXE_PREFIX="\$APPDIR/$GTK_EXEC_PREFIX"
export GTK_PATH="\$APPDIR/usr/lib/gtk-2.0"
export GTK_IM_MODULE_FILE="\$APPDIR/usr/lib/gtk-2.0/$GTK_BINARY_VERSION/immodules.cache"
EOF
        ;;
    3)
        echo "Processing GTK 3 modules..."
        GTK_LIBDIR="$(pkg_var libdir gtk+-3.0 "$LIB_PATH")/gtk-3.0"
        GTK_BINARY_VERSION="$(pkg_var gtk_binary_version gtk+-3.0 "3.0.0")"

        # 复制 immodules
        IMMODULES_DIR="$GTK_LIBDIR/$GTK_BINARY_VERSION/immodules"
        if [ -d "$IMMODULES_DIR" ]; then
            mkdir -p "$APPDIR/usr/lib/gtk-3.0/$GTK_BINARY_VERSION/immodules"
            cp -a "$IMMODULES_DIR"/*.so "$APPDIR/usr/lib/gtk-3.0/$GTK_BINARY_VERSION/immodules/" 2>/dev/null || true
            if command -v gtk-query-immodules-3.0 > /dev/null; then
                gtk-query-immodules-3.0 > "$APPDIR/usr/lib/gtk-3.0/$GTK_BINARY_VERSION/immodules.cache" 2>/dev/null || true
            fi
        fi

        # 复制 printbackends
        PRINTBACKENDS_DIR="$GTK_LIBDIR/$GTK_BINARY_VERSION/printbackends"
        if [ -d "$PRINTBACKENDS_DIR" ]; then
            mkdir -p "$APPDIR/usr/lib/gtk-3.0/$GTK_BINARY_VERSION/printbackends"
            cp -a "$PRINTBACKENDS_DIR"/*.so "$APPDIR/usr/lib/gtk-3.0/$GTK_BINARY_VERSION/printbackends/" 2>/dev/null || true
        fi

        GTK_EXEC_PREFIX="$(pkg_var exec_prefix gtk+-3.0 /usr)"
        cat >> "$HOOK_FILE" << EOF
export GTK_EXE_PREFIX="\$APPDIR/$GTK_EXEC_PREFIX"
export GTK_PATH="\$APPDIR/usr/lib/gtk-3.0"
export GTK_IM_MODULE_FILE="\$APPDIR/usr/lib/gtk-3.0/$GTK_BINARY_VERSION/immodules.cache"
EOF
        ;;
esac

# 复制 GDK PixBuf
echo "Copying GDK PixBuf..."
GDK_LIBDIR="$(pkg_var libdir gdk-pixbuf-2.0 "$LIB_PATH")"
GDK_BINARY_DIR="$(pkg_var gdk_pixbuf_binarydir gdk-pixbuf-2.0 "$GDK_LIBDIR/gdk-pixbuf-2.0/2.10.0")"
GDK_MODULE_DIR="$(pkg_var gdk_pixbuf_moduledir gdk-pixbuf-2.0 "$GDK_BINARY_DIR/loaders")"
GDK_CACHE_FILE="$(pkg_var gdk_pixbuf_cache_file gdk-pixbuf-2.0 "$GDK_BINARY_DIR/loaders.cache")"

if [ -d "$GDK_BINARY_DIR" ]; then
    mkdir -p "$APPDIR/usr/lib/gdk-pixbuf-2.0/2.10.0/loaders"
    cp -a "$GDK_MODULE_DIR"/*.so "$APPDIR/usr/lib/gdk-pixbuf-2.0/2.10.0/loaders/" 2>/dev/null || true
    # 更新 cache
    if command -v gdk-pixbuf-query-loaders > /dev/null; then
        gdk-pixbuf-query-loaders > "$APPDIR/usr/lib/gdk-pixbuf-2.0/2.10.0/loaders.cache" 2>/dev/null || true
    fi
    echo "export GDK_PIXBUF_MODULE_FILE=\"\$APPDIR/usr/lib/gdk-pixbuf-2.0/2.10.0/loaders.cache\"" >> "$HOOK_FILE"
fi

# 复制核心库
echo "Copying core libraries..."
LIBS=(
    "gdk-pixbuf-2.0:libgdk_pixbuf-*.so*"
    "gobject-2.0:libgobject-*.so*"
    "gio-2.0:libgio-*.so*"
    "librsvg-2.0:librsvg-*.so*"
    "pango:libpango-*.so*"
    "pangocairo:libpangocairo-*.so*"
    "pangoft2:libpangoft2-*.so*"
)

LINUXDEPLOY_LIBS=()
for entry in "${LIBS[@]}"; do
    lib="${entry%%:*}"
    pattern="${entry##*:}"
    dir="$(pkg_var libdir "$lib" "$LIB_PATH")"
    if [ -d "$dir" ]; then
        while IFS= read -r -d '' file; do
            LINUXDEPLOY_LIBS+=("--library=$file")
        done < <(find "$dir" \( -type l -o -type f \) -name "$pattern" -print0 2>/dev/null)
    fi
done

# 调用 linuxdeploy 处理库
if [ -n "$LINUXDEPLOY" ] && [ ${#LINUXDEPLOY_LIBS[@]} -gt 0 ]; then
    echo "Running linuxdeploy for library processing..."
    env LINUXDEPLOY_PLUGIN_MODE=1 "$LINUXDEPLOY" --appdir="$APPDIR" "${LINUXDEPLOY_LIBS[@]}"
fi

# 处理 WebKit 依赖
if [ -n "$WEBKIT_VERSION" ]; then
    echo "Processing WebKit $WEBKIT_VERSION..."
    WEBKIT_LIBS=()

    # 根据 WebKit 版本查找库
    case "$WEBKIT_VERSION" in
        4.0)
            WEBKIT_PATTERN="libwebkit2gtk-4.0*.so*"
            ;;
        4.1)
            WEBKIT_PATTERN="libwebkit2gtk-4.1*.so*"
            ;;
        *)
            echo "WARNING: Unknown WebKit version: $WEBKIT_VERSION"
            WEBKIT_PATTERN=""
            ;;
    esac

    if [ -n "$WEBKIT_PATTERN" ]; then
        # 查找 WebKit 库
        while IFS= read -r -d '' file; do
            WEBKIT_LIBS+=("--library=$file")
        done < <(find "$LIB_PATH" \( -type l -o -type f \) -name "$WEBKIT_PATTERN" -print0 2>/dev/null)

        # 查找 WebKit 辅助进程
        for helper in WebKitWebProcess WebKitNetworkProcess; do
            while IFS= read -r -d '' file; do
                # 复制到 usr/lib/webkit2gtk-4.0 或 webkit2gtk-4.1
                rel_dir="usr/lib/webkit2gtk-${WEBKIT_VERSION}"
                mkdir -p "$APPDIR/$rel_dir"
                cp -a "$file" "$APPDIR/$rel_dir/"
            done < <(find /usr/lib -name "$helper" -type f -print0 2>/dev/null)
        done

        # 复制 injected bundle
        BUNDLE_NAME="libwebkit2gtkinjectedbundle.so"
        while IFS= read -r -d '' file; do
            rel_dir="usr/lib/webkit2gtk-${WEBKIT_VERSION}"
            mkdir -p "$APPDIR/$rel_dir"
            cp -a "$file" "$APPDIR/$rel_dir/"
        done < <(find /usr/lib -name "$BUNDLE_NAME" -type f -print0 2>/dev/null)

        # 调用 linuxdeploy 处理 WebKit 库
        if [ -n "$LINUXDEPLOY" ] && [ ${#WEBKIT_LIBS[@]} -gt 0 ]; then
            echo "Running linuxdeploy for WebKit libraries..."
            env LINUXDEPLOY_PLUGIN_MODE=1 "$LINUXDEPLOY" --appdir="$APPDIR" "${WEBKIT_LIBS[@]}"
        fi
    fi
fi

chmod +x "$HOOK_FILE"
echo "GTK bundle completed."
