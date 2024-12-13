package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	// "path/filepath"
	// "github.com/nfnt/resize"
)

// compressImage 压缩图片到指定大小
func compressImage(inputPath, outputPath string, maxSizeKB int) error {
	// 打开图片
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// 计算当前图片的大小（以字节为单位）
	currentSize, err := fileSize(inputPath)
	if err != nil {
		return fmt.Errorf("failed to get file size: %v", err)
	}

	// 如果当前大小已经小于或等于最大大小，则不需要压缩
	if currentSize <= int64(maxSizeKB*1024) {
		fmt.Println("The image is already smaller than or equal to the maximum size.")
		return saveImage(outputPath, img, format, 95)
	}

	// 尝试不同的质量值，直到图片大小小于或等于最大大小
	quality := 95 // 从95%的质量开始尝试
	for quality > 0 {
		// 保存图片到临时路径
		tempPath := outputPath + "_temp"
		if format == "jpeg" || format == "jpg" {
			tempPath += ".jpg"
		} else if format == "png" {
			tempPath += ".png"
		}
		err = saveImage(tempPath, img, format, quality)
		if err != nil {
			return fmt.Errorf("failed to save image: %v", err)
		}

		// 检查临时图片的大小
		newSize, err := fileSize(tempPath)
		if err != nil {
			return fmt.Errorf("failed to get file size: %v", err)
		}

		// 如果大小小于或等于最大大小，则将临时图片重命名为最终输出路径
		if newSize <= int64(maxSizeKB*1024) {
			fmt.Printf("Image compressed to %d bytes with quality %d.\n", newSize, quality)
			return os.Rename(tempPath, outputPath)
		}

		fmt.Printf("Now, Image size: %d, quality: %d\n", newSize, quality)

		quality -= 5
	}

	return fmt.Errorf("failed to compress image to the desired size")
}

// fileSize 获取文件大小
func fileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// saveImage 保存图片到指定路径
func saveImage(path string, img image.Image, format string, quality int) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if format == "jpeg" || format == "jpg" {
		options := &jpeg.Options{Quality: quality}
		return jpeg.Encode(out, img, options)
	} else if format == "png" {
		encoder := png.Encoder{CompressionLevel: png.CompressionLevel(9 - quality/10)}
		return encoder.Encode(out, img)
	}
	return fmt.Errorf("unsupported image format: %s", format)
}

func main() {
	inputPath := "input.jpg"
	outputPath := "output.jpg"
	maxSizeKB := 500

	err := compressImage(inputPath, outputPath, maxSizeKB)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
