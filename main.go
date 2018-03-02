package main

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"os"
)

const MaxDisplaySize int = 1000000 //1MB
const MaxDisplayWidth int = 1024   //1024pixel
const MaxDisplayHeight int = 1024  //1024pixel

const MaxLimitSize int = 6000000 //6MB, AWS Api Gateway Max Response Body Size (2018.3)

func main() {
	// open "test.jpg"
	file, err := os.Open("test.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(1000, 0, img, resize.Lanczos3)

	out, err := os.Create("test_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}

//Resize when image bigger than max display size
//1.resize width/height if bigger than max display width/height
//2.compress quility
//3.if resized image is bigger than max limit size!? give up! Orz
func resize() {
}
