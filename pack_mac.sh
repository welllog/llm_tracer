#!/bin/bash
set -e

APP_NAME="LLMTracer"
APP_DIR="${APP_NAME}.app"
CONTENTS_DIR="${APP_DIR}/Contents"
MAC_OS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"

echo "=== 1. 一键编译前端与后端 ==="
make build

echo "=== 2. 创建 macOS App 目录结构 ==="
rm -rf "${APP_DIR}"
mkdir -p "${MAC_OS_DIR}"
mkdir -p "${RESOURCES_DIR}"

echo "=== 3. 搬运可执行二进制 ==="
cp llm_tracer_web "${MAC_OS_DIR}/${APP_NAME}"
chmod +x "${MAC_OS_DIR}/${APP_NAME}"

echo "=== 4. 生成 App 图标 (.icns) ==="
if [ -f "logo.png" ]; then
    echo "检测到 logo.png，正在转换为 App 图标..."
    mkdir -p temp_icon.iconset
    
    # 转换为各个尺寸的 png 图标
    sips -s format png -z 16 16     logo.png --out temp_icon.iconset/icon_16x16.png
    sips -s format png -z 32 32     logo.png --out temp_icon.iconset/icon_16x16@2x.png
    sips -s format png -z 32 32     logo.png --out temp_icon.iconset/icon_32x32.png
    sips -s format png -z 64 64     logo.png --out temp_icon.iconset/icon_32x32@2x.png
    sips -s format png -z 128 128   logo.png --out temp_icon.iconset/icon_128x128.png
    sips -s format png -z 256 256   logo.png --out temp_icon.iconset/icon_128x128@2x.png
    sips -s format png -z 256 256   logo.png --out temp_icon.iconset/icon_256x256.png
    sips -s format png -z 512 512   logo.png --out temp_icon.iconset/icon_256x256@2x.png
    sips -s format png -z 512 512   logo.png --out temp_icon.iconset/icon_512x512.png
    sips -s format png -z 1024 1024 logo.png --out temp_icon.iconset/icon_512x512@2x.png
    
    # 利用系统工具合并为 icns
    iconutil -c icns temp_icon.iconset
    mv temp_icon.icns "${RESOURCES_DIR}/icon.icns"
    rm -rf temp_icon.iconset
    echo "图标生成成功！"
else
    echo "⚠️ 未找到 logo.png，将无法为 App 设置图标。"
fi

echo "=== 5. 生成 Info.plist 配置文件 ==="
cat > "${CONTENTS_DIR}/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>${APP_NAME}</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.ludao.llmtracer</string>
    <key>CFBundleName</key>
    <string>${APP_NAME}</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>LSUIElement</key>
    <string>0</string>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF

echo "=== 🎉 macOS App 打包完成！ ==="
echo "生成路径：$(pwd)/${APP_DIR}"
echo "现在您可以双击 ${APP_DIR} 直接静默启动运行！"
