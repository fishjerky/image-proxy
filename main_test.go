package main

import (
	"testing"
	//"github.com/disintegration/imaging"
)

//1. Normal case
// png/jpg/gif
func Test_Normal(t *testing.T) {
	url := "https://upload.wikimedia.org/wikipedia/commons/thumb/4/47/PNG_transparency_demonstration_1.png/280px-PNG_transparency_demonstration_1.png"

	image, err := GetImageFromUrl(url)
	if err != nil {
		t.Error("就是不通過")
	}
	t.Log("[case]Normal pass, size: %d", len(image))
}

//2.resize image
// -3m/6m/14m
func Test_Resize(t *testing.T) {
	/*	t.Skip("skipping this case")
		filePath := "testdata/10m.jpg"
		src, err := imaging.Open(filePath)
		if err != nil {
			t.Error("Open failed: %v", err)
		}

		resized := resize(src)
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
	referer := "http://myfone.taiwanmobile.com/"
	if CheckReferer(referer) == false {
		t.Error("Error! " + referer + " should be a valid referer")
	}

	referer = "http://myfone.taiwanmobile.com/"
	if CheckReferer(referer) == false {
		t.Error("Error! " + referer + " should be a valid referer")
	}

	//invalid
	referer = "http://www.pchome.com/"
	if CheckReferer(referer) == false {
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
