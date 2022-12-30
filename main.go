package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/image/draw"
)

const IMG_DIR = "/mnt/c/Users/hiros/Documents/画面キャプチャ/"
const DEST_DIR = "./images/"
const BASE_WIDTH = 1920
const BASE_HEIGHT = 1080

func saveImage(baseFilename string, postFix string, srcImage image.Image, srcRect image.Rectangle) {
	var newFilename string
	if postFix == "" {
		newFilename = filepath.Join(DEST_DIR, baseFilename+".jpg")
	} else {
		newFilename = filepath.Join(DEST_DIR, baseFilename+"-"+postFix+".jpg")
	}
	dstF, err := os.Create(newFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer dstF.Close()
	dstWidth := srcRect.Size().X
	dstHeight := srcRect.Size().Y
	if dstWidth > BASE_WIDTH || dstHeight > BASE_HEIGHT {
		dstHeight = dstHeight * BASE_WIDTH / dstWidth
		dstWidth = BASE_WIDTH
	}
	fmt.Println(srcImage.Bounds(), srcRect.Size().X, srcRect.Size().Y, dstWidth, dstHeight)
	newImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{dstWidth, dstHeight}})
	draw.BiLinear.Scale(newImage, newImage.Bounds(), srcImage, srcRect, draw.Over, nil)
	jpeg.Encode(dstF, newImage, &jpeg.Options{Quality: 90})
}

func processImage(filename string, modTime time.Time) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	srcImage, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	baseFilename := modTime.Format("20060102-150405")
	width := srcImage.Bounds().Size().X
	height := srcImage.Bounds().Size().Y
	fmt.Println(baseFilename, width, height)
	if width == 5760 && height == 2160 {
		saveImage(baseFilename, "01", srcImage, image.Rect(0, 0, 3840, 1920))
		saveImage(baseFilename, "02", srcImage, image.Rect(3840, 540, 5760, 1620))
	} else {
		saveImage(baseFilename, "00", srcImage, srcImage.Bounds())
	}
}

func main() {
	files, err := ioutil.ReadDir(IMG_DIR)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fn := filepath.Join(IMG_DIR, file.Name())
		fmt.Println(fn, file.ModTime())
		processImage(fn, file.ModTime())
	}
	fmt.Println("Convert imgs")
}
