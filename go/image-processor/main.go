package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: image-processor <input_image> <output_image>")
		return
	}

	inputPath := os.Args[1]
	outPath := os.Args[2]

	// Open the image
	img, err := openImage(inputPath)
	if err != nil {
		fmt.Printf("Error opening image: %v\n", err)
		return
	}

	fmt.Printf("Processing image: %s\n", inputPath)

	// Convert to grayscale
	grayImg := grayscale(img)

	// Save the grayscale image
	err = saveImage(grayImg, outPath)
	if err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		return
	}

	fmt.Printf("Grayscale image saved to: %s\n", outPath)
}

// openImage opens and decodes an image from the given path.
func openImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// saveImage encodes and saves an image to the given path.
func saveImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	ext := filepath.Ext(path)

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(f, img, nil)
	case ".png":
		return png.Encode(f, img)
	default:
		return fmt.Errorf("unsupported output format: %s", ext)
	}
}

// grayscale converts an image to grayscale.
func grayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	return gray
}