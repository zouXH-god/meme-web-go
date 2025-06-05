package memes_cli

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
)

// CombineImagesVertically 垂直拼接多张图片，先统一宽度再拼接
// 支持JPEG和PNG格式，输出为PNG格式
func CombineImagesVertically(imagesData [][]byte) (string, error) {
	// 1. 解码所有图片并找出最大宽度
	var images []image.Image
	maxWidth := 0
	var totalHeight int

	for _, imgData := range imagesData {
		img, _, err := image.Decode(bytes.NewReader(imgData))
		if err != nil {
			return "", err
		}

		bounds := img.Bounds()
		if bounds.Dx() > maxWidth {
			maxWidth = bounds.Dx()
		}

		images = append(images, img)
	}

	// 2. 预处理所有图片：统一宽度并计算总高度
	var processedImages []image.Image
	totalHeight = 0

	for _, img := range images {
		bounds := img.Bounds()
		imgWidth := bounds.Dx()
		imgHeight := bounds.Dy()

		// 计算缩放后的高度（保持宽高比）
		newHeight := int(math.Round(float64(imgHeight) * float64(maxWidth) / float64(imgWidth)))

		// 创建缩放后的图像
		scaledImg := image.NewRGBA(image.Rect(0, 0, maxWidth, newHeight))

		// 使用双线性插值进行缩放（比简单缩放质量更好）
		for y := 0; y < newHeight; y++ {
			srcY := float64(y) * float64(imgHeight) / float64(newHeight)
			y1 := int(math.Floor(srcY))
			y2 := int(math.Ceil(srcY))
			if y2 >= imgHeight {
				y2 = imgHeight - 1
			}

			for x := 0; x < maxWidth; x++ {
				srcX := float64(x) * float64(imgWidth) / float64(maxWidth)
				x1 := int(math.Floor(srcX))
				x2 := int(math.Ceil(srcX))
				if x2 >= imgWidth {
					x2 = imgWidth - 1
				}

				// 双线性插值
				c1 := interpolateColor(img.At(x1, y1), img.At(x2, y1), srcX-float64(x1))
				c2 := interpolateColor(img.At(x1, y2), img.At(x2, y2), srcX-float64(x1))
				finalColor := interpolateColor(c1, c2, srcY-float64(y1))

				scaledImg.Set(x, y, finalColor)
			}
		}

		processedImages = append(processedImages, scaledImg)
		totalHeight += newHeight
	}

	// 3. 创建最终图片并拼接
	newImg := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))
	currentY := 0

	for _, img := range processedImages {
		bounds := img.Bounds()
		draw.Draw(newImg,
			image.Rect(0, currentY, bounds.Dx(), currentY+bounds.Dy()),
			img,
			image.Point{},
			draw.Src)
		currentY += bounds.Dy()
	}

	// 4. 编码为PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, newImg); err != nil {
		return "", err
	}

	return SaveLocalImage(buf.Bytes())
}

// interpolateColor 颜色插值函数
func interpolateColor(c1, c2 color.Color, ratio float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return color.RGBA{
		R: uint8(float64(r1>>8)*(1-ratio) + float64(r2>>8)*ratio),
		G: uint8(float64(g1>>8)*(1-ratio) + float64(g2>>8)*ratio),
		B: uint8(float64(b1>>8)*(1-ratio) + float64(b2>>8)*ratio),
		A: uint8(float64(a1>>8)*(1-ratio) + float64(a2>>8)*ratio),
	}
}
