package main

import (
	"image"
	_ "image/color"
	"os"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/disintegration/imaging"
)

const MaxDisplaySize int = 1000000 //1MB
const MaxDisplayWidth int = 1024   //1024pixel
const MaxDisplayHeight int = 1024  //1024pixel
const MaxLimitSize int = 6000000   //6MB, AWS Api Gateway Max Response Body Size (2018.3)

type Response struct {
	statusCode int
}

func main() {
	log.SetOutput(os.Stdout)

	//1.get image
	url := ""
	GetImageFromUrl(url)

	//2.resize
	filePath := "testdata/14m.jpg"
	width, height := getImageDimensionFromFile(filePath)
	log.Printf("Width:%d, Height: %d", width, height)

	// open "test.jpg"
	src, err := imaging.Open(filePath)
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}

	resize(src, width, height)

	//3.response
}

func GetImageFromUrl(imageUrl string) ([]byte, error) {

	var body []byte
	resp, err := http.Get(imageUrl)
	if err != nil {
		log.Fatal(err)
	} else {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Code: %d Url: %s", resp.StatusCode, imageUrl)

	//exception case handling

	// 304是讀cache，res.headers.location可以cover這個情境, 所以不用判斷
	if resp.StatusCode >= 300 && resp.StatusCode <= 400 && resp.Header.Get("location") != "" {
		// && res.headers.location
		log.Printf("\t- REDIRECT DECTECTED!!! Status Code: %d", resp.StatusCode)

		// The location for some (most) redirects will only contain the path,  not the hostname;
		// detect this and add the host to the path.
		u, err := url.Parse(resp.Header.Get("location"))
		if err != nil {
			log.Fatal(err)
		}

		if u.Host != "" {
			log.Printf("hostname: %s", u.Host)
			// Hostname included; make request to res.headers.location
			redirectUrl := resp.Header.Get("location")
			log.Printf("redirectUrl = %+v\n", redirectUrl)
		} else {
			log.Fatalf("Missing hostname: %s", resp.Header.Get("location"))
			//var myError = new Error(msg, "Missing hostname");
			//successCallback(myError);
		}
	}

	return body, err
}

//Resize when image bigger than max display size
//1.resize width/height if bigger than max display width/height
//2.compress quility
//3.if resized image is bigger than max limit size!? give up! Orz
func resize(src image.Image, width int, height int) {
	//1. resize if bigger than max display width/height
	switch {
	case width > MaxDisplayWidth:
		src = imaging.Resize(src, MaxDisplayWidth, 0, imaging.NearestNeighbor)
	case height > MaxDisplayHeight:
		src = imaging.Resize(src, 0, MaxDisplayHeight, imaging.NearestNeighbor)
	}

	//	dst := imaging.New(512, 512, color.NRGBA{0, 0, 0, 0})

	// Save the resulting image using JPEG format.
	err := imaging.Save(src, "result/out_example.jpg")
	if err != nil {
		log.Fatalf("Save failed: %v", err)
	}
}

func getImageDimensionFromFile(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("%v", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
