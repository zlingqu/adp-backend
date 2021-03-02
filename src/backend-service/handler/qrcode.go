package handler

import (
	"github.com/gin-gonic/gin"
	qrcode "github.com/skip2/go-qrcode"
	"net/http"
)

type ReqQrcode struct {
	Url string `form:"url"`
}

func GetQrcode(c *gin.Context) {
	var reqQrcode ReqQrcode
	// var code int
	var png []byte
	c.ShouldBind(&reqQrcode)

	png, _ = qrcode.Encode(reqQrcode.Url, qrcode.Medium, 256)
	// png, _ = json.Marshal(png)
	c.String(http.StatusOK, string(png))
	// c.Data(http.StatusOK, "bytes", png)
}