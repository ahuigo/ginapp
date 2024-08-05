package tpls

import (
	"ginapp/utils/String"
	"io"
	"net/http"
	URL "net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// requestForm('post','http://m:4500/tpl/redirect?redirect_uri=https://s/dump/ab/c/c2',{redirect_uri:"https://s/dump/ab/c/c"})
// or: curl 'http://m:4500/tpl/redirect?redirect_uri=https://s/dump/ab/c/c2' -d '{id_token:"xx"}'
func TplRedirectPage(ctx *gin.Context) {
	url := String.EncodeURI(ctx.Query("redirect_uri"))
	if url == "" {
		url = "https://s/dump/ab/c/c2"
	}
	body, _ := io.ReadAll(ctx.Request.Body)
	query, _ := URL.ParseQuery(string(body))
	data := map[string]string{}
	for k, v := range query {
		data[k] = v[0]
	}
	rdata := map[string]interface{}{
		"url":  url,
		"now1": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
		"body": data,
	}

	_ = rdata
	// ctx.JSON(http.StatusOK, "redirect.tmpl")
	ctx.HTML(http.StatusOK, "redirect.tmpl", rdata)

}
