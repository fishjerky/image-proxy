package main

import (
	"bytes"
	"errors"
	"image"
	_ "image/color"
	"image/png"
	"os"
	"strings"

	"encoding/base64"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/disintegration/imaging"
	//github.com/nfnt/resize //no longer being updated
	"io/ioutil"
	"log"
	"net/http"
	//"net/url"
)

const MaxDisplaySize int = 1000000 //1MB
const MaxDisplayWidth int = 1024   //1024pixel
const MaxDisplayHeight int = 1024  //1024pixel
const MaxLimitSize int = 6000000   //6MB, AWS Api Gateway Max Response Body Size (2018.3)

type Response struct {
	statusCode int
}

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrPicUrlNotProvided = errors.New("no pic url was provided in the URL")
	ErrInvalidReferer    = errors.New("Invalid api call")
	ErrResizeFailed      = errors.New("Resizing image failed")
	ErrFetchImageFailed  = errors.New("Fetching image failed")
	ErrDecodeImageFailed = errors.New("Decoding image failed")
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	imgUrl := request.QueryStringParameters["p"]
	referer := request.Headers["Referer"]
	//log.Printf("QueryStringParameters : %+v", request.QueryStringParameters)
	//log.Printf("Headers : %+v", request.Headers)
	log.Printf("Image url: %s, Referer: %s", imgUrl, referer)

	// If no name is provided in the HTTP request body, throw an error
	/*
		if len(request.QueryStringParameters.p) < 1 {
			return events.APIGatewayProxyResponse{}, ErrPicUrlNotProvided
		}
	*/

	if !CheckReferer(referer) {
		return events.APIGatewayProxyResponse{}, ErrInvalidReferer
	}

	//fetch image
	byteImg, err := GetImageFromUrl(imgUrl)
	if err != nil {
		switch err {
		case ErrPicUrlNotProvided:
			return events.APIGatewayProxyResponse{}, err
		default:
			log.Fatal(err)
			return events.APIGatewayProxyResponse{}, ErrFetchImageFailed
		}
	}

	//resize image, if needed
	byteResized, err2 := resize(byteImg)
	if err2 != nil {
		log.Fatal(err)
		return events.APIGatewayProxyResponse{}, ErrResizeFailed
	}
	/*
		//prepare
		// create buffer
		buff := new(bytes.Buffer)

		// encode image to buffer
		err = png.Encode(buff, imgImg)
		if err != nil {
			log.Fatal("failed to create buffer", err)
		}
		// convert buffer to reader
		dist := make([]byte, len(buff.Bytes())*3) //base64 3 times bigger
	*/
	dist := make([]byte, len(byteResized)*3) //base64 3 times bigger
	base64.StdEncoding.Encode(dist, byteResized)

	return events.APIGatewayProxyResponse{
		Body:       string(dist),
		StatusCode: 200,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(Handler)

}

func GetImageFromUrl(imageUrl string) ([]byte, error) {
	if imageUrl == "" {
		return nil, ErrPicUrlNotProvided
	}

	var body []byte
	resp, err := http.Get(imageUrl)
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	//exception case handling
	//handleRedirect(resp )  //go seems handle it already

	return body, err
}

/*
func handleRedirect(resp HttpResponse) {
	// 304是讀cache，res.headers.location可以cover這個情境, 所以不用判斷
	if resp.StatusCode >= 300 && resp.StatusCode <= 400 && resp.Header.Get("location") != "" {
		// && res.headers.location
		log.Printf("\t- REDIRECT DECTECTED!!! Status Code: %d", resp.StatusCode)

		// The location for some (most) redirects will only contain the path,  not the hostname;
		// detect this and add the host to the path.
		u, err := url.Parse(resp.Header.Get("location"))
		if err != nil {
			log.Fatal(err)
			return nil, err
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

}
*/

//Resize when image bigger than max display size
//1.resize width/height if bigger than max display width/height
//2.compress quility
//3.if resized image is bigger than max limit size!? give up! Orz
func resize(byteImg []byte) ([]byte, error) {
	//1. resize if bigger than max display width/height
	orgImg, _, err := image.Decode(bytes.NewReader(byteImg))
	if err != nil {
		log.Fatal(err)
		return nil, ErrDecodeImageFailed
	}

	width := orgImg.Bounds().Dx()
	height := orgImg.Bounds().Dy()

	//prevent empty
	if (width == 0) && (height == 0) {
		return byteImg, nil
	}

	if width < MaxDisplayWidth && height < MaxDisplayHeight {
		return byteImg, nil
	}

	log.Printf("Picture is too big(w:%d, h:%d, size:%d). Start to resize image", width, height, len(byteImg))
	//resize big image
	resized := orgImg
	switch {
	case width > MaxDisplayWidth:
		resized = imaging.Resize(orgImg, MaxDisplayWidth, 0, imaging.NearestNeighbor)
	case height > MaxDisplayHeight:
		resized = imaging.Resize(orgImg, 0, MaxDisplayHeight, imaging.NearestNeighbor)
	}

	buff := new(bytes.Buffer)

	// encode image to buffer
	err = png.Encode(buff, resized)
	if err != nil {
		log.Fatal("failed to create buffer", err)
		return nil, err
	}

	// Save the resulting image using JPEG format.
	/*err = imaging.Save(orgImg, "temp/org.png")
	//	err = imaging.Save(resized, "temp/re.png")
	if err != nil {
		log.Fatalf("Save failed: %v", err)
	}*/
	log.Printf("Finish resizing from %d(w:%d, h:%d) to %d(w:%d, h:%d)",
		len(byteImg), width, height, len(buff.Bytes()), resized.Bounds().Dx(), resized.Bounds().Dy())

	return buff.Bytes(), nil
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

func CheckReferer(referer string) bool {
	validReferers := []string{"myfone.com.tw", "taiwanmobile.com"}
	//referer = request.Header.Get("Referer")
	if len(referer) == 0 {
		log.Printf("referer is empty")
		return false
	}

	for _, domain := range validReferers {
		if strings.Contains(referer, domain) == true {
			return true
		}

	}
	return false
}
