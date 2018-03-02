package main

import (
	"testing"
)

//1. Normal case
// png/jpg/gif
func Test_Normal(t *testing.T) {
	//		t.Error("除法函數測試沒通過") //如果不是如預期的那麼就報錯
	//		t.Log("第一個測試通過了") //記錄一些你期望記錄的信息
}

//2.resize image
// -3m/6m/14m
func Test_Resize(t *testing.T) {
	t.Error("就是不通過")
}

//3.boundary case
// ->=<1m
//	-1m https://upload.wikimedia.org/wikipedia/commons/f/fc/LythrumSalicaria-flower-1mb.jpg
// ->+<6m
//	-14m http://farm4.static.flickr.com/3182/2893346171_11a4df8533_o.jpg

func Test_Boundary(t *testing.T) {
}

func Test_Redirect(t *testing.T) {
}

//Test case
//4.redirect case
// http://www.myfone.com.tw/buy/myfoneweb/buy/Download_app/button.png
// domain redirect
//http://myfone.taiwanmobile.com/buy/myfoneweb/buy/Download_app/button.png
//5.Referer
// myfone.taiwanmobile.com
//6.url encode
// http://www.0958566678.url.tw/NEW/NEW%2520Phone/Y7max.png
