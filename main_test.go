package main

import (
	//"github.com/disintegration/imaging"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var (

//query["p"] = "http://files.softicons.com/download/game-icons/super-mario-icons-by-sandro-pereira/png/16/Mushroom%20-%20Super.png"
)

//1. Normal case
// png/jpg/gif
func Test_Normal(t *testing.T) {
	url := "https://upload.wikimedia.org/wikipedia/commons/thumb/4/47/PNG_transparency_demonstration_1.png/280px-PNG_transparency_demonstration_1.png"

	image, err := GetImageFromUrl(url)
	if err != nil {
		t.Error("就是不通過")
	}
	t.Log("[case]Normal pass, size: ", len(image))
}

//2.resize image
// -3m/6m/14m
func Test_Resize(t *testing.T) {
	/*
		filePath := "testdata/10m.jpg"
		src, err := os.Open(filePath)
		if err != nil {
			t.Error("Open failed: %v", err)
		}

		resized := resize(src)

		if resized.Bounds().Dx() > MaxDisplayWidth {
			t.Error("Resize failed: ")
		}
	*/
}

//3.boundary case
// ->=<1m
//	-1m https://upload.wikimedia.org/wikipedia/commons/f/fc/LythrumSalicaria-flower-1mb.jpg
// ->+<6m
//	-14m http://farm4.static.flickr.com/3182/2893346171_11a4df8533_o.jpg

func Test_Boundary(t *testing.T) {
	//t.Error("就是不通過")
}

//4.redirect case
func Test_Redirect(t *testing.T) {

	//normail
	url := "http://www.myfone.com.tw/buy/myfoneweb/buy/Download_app/button.png"
	image, err := GetImageFromUrl(url)
	if err != nil {
		t.Error("Error")
	}
	t.Log("[Case] Redirect pass with size: ", len(image))

	//domain
	url = "http://myfone.taiwanmobile.com/buy/myfoneweb/buy/Download_app/button.png"
	image, err = GetImageFromUrl(url)
	if err != nil {
		t.Error("Error")
	}
	t.Log("[case] Redirect domain pass with size: ", len(image)) //記錄一些你期望記錄的信息

}

//Test case
//5.Referer
func Test_CheckReferer(t *testing.T) {
	//valid
	referer := "http://myfone.com.tw/"
	if !CheckReferer(referer) {
		t.Error("Error! " + referer + " should be a valid referer")
	}

	referer = "http://myfone.taiwanmobile.com/"
	if !CheckReferer(referer) {
		t.Error("Error! " + referer + " should be a valid referer")
	}

	//invalid
	referer = "http://www.pchome.com/"
	if CheckReferer(referer) {
		t.Error("Error! " + referer + " is not a valid referer")
	}

}

//6.url encode
func Test_UrlEncode(t *testing.T) {
	url := "http://www.0958566678.url.tw/NEW/NEW%2520Phone/Y7max.png"
	image, err := GetImageFromUrl(url)
	if err != nil {
		t.Error("Error")
	}
	t.Log("[Case] Redirect pass with size: ", len(image))

}

//7. http/https
/*
func TestHandler(t *testing.T) {
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: "Paul"},
			expect:  "Hello Paul",
			err:     nil,
		},
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: ""},
			expect:  "",
			err:     ErrNameNotProvided,
		},
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}*/
func TestHandler(t *testing.T) {
	headers := make(map[string]string)
	query1 := make(map[string]string)
	headers["Referer"] = "www.myfone.com.tw"
	query1["p"] = "http://files.softicons.com/download/game-icons/super-mario-icons-by-sandro-pereira/png/16/Mushroom%20-%20Super.png"

	query2 := make(map[string]string)
	query2["p"] = "https://vignette.wikia.nocookie.net/fantendo/images/6/6e/Small-mario.png/revision/latest?cb=20120718024112"
	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  int
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Headers: headers, QueryStringParameters: query1},
			expect:  200,
			err:     nil,
		},
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Headers: headers, QueryStringParameters: query2},
			expect:  200,
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.StatusCode)
	}
}
func TestInvalidHandler(t *testing.T) {
	headers := make(map[string]string)
	query := make(map[string]string)
	headers["Referer"] = "invalid.referer.com"
	//query["p"] = "http://files.softicons.com/download/game-icons/super-mario-icons-by-sandro-pereira/png/16/Mushroom%20-%20Super.png"

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Headers: headers, QueryStringParameters: query},
			expect:  "",
			err:     ErrPicUrlNotProvided,
		},
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Headers: headers},
			expect:  "",
			err:     ErrInvalidReferer,
		},
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
