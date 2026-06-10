import os
from PIL import Image, ImageDraw

def main():
    icon_path = "logo.png"
    if not os.path.exists(icon_path):
        print(f"Error: {icon_path} not found.")
        return

    # 1. 打开原图
    img = Image.open(icon_path).convert("RGBA")
    w, h = img.size
    print(f"Original image size: {w}x{h}")

    # 2. 创建一个同尺寸的全透明背景图像
    output = Image.new("RGBA", (w, h), (0, 0, 0, 0))

    # 3. 构造 macOS Squircle 圆角矩形遮罩
    # 图标在 1024x1024 像素的图像中占据中间 [164, 164, 860, 860] 的区域
    left = 164
    top = 164
    right = 860
    bottom = 860
    box_w = right - left
    radius = int(box_w * 0.223) # macOS 标准圆角比例约 22.3%

    # 4. 在 Mask 上绘制白色的圆角矩形
    mask = Image.new("L", (w, h), 0)
    mask_draw = ImageDraw.Draw(mask)
    mask_draw.rounded_rectangle(
        [left, top, right, bottom],
        radius=radius,
        fill=255
    )

    # 5. 复合图像：将原图绘制到透明背景上，并且只保留遮罩范围内的像素
    output.paste(img, (0, 0), mask)

    # 6. 保存覆盖原图
    output.save(icon_path, "PNG")
    print(f"App icon masked successfully with transparent background!")

if __name__ == "__main__":
    main()
