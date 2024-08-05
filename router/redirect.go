package router

import (
	"ginapp/utils/String"
	"io"
	"net/http"
	URL "net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// EchoHandler _
func RedirectServer(ctx *gin.Context) {
	redirectUri := ctx.Query("redirect_uri")
	switch ctx.Param("code") {
	case "302":
		ctx.Redirect(http.StatusFound, redirectUri)
	case "307":
		ctx.Redirect(http.StatusTemporaryRedirect, redirectUri)
	case "308":
		ctx.Redirect(http.StatusPermanentRedirect, redirectUri)
	case "refresh":
		redirectRefresh(ctx)
	case "form":
		fallthrough
	default:
		redirectForm(ctx)
	}
}

// RedirectHandler
func redirectRefresh(ctx *gin.Context) {
	url := String.EncodeURI(ctx.Query("redirect_uri"))
	html := `<head>
	<meta http-equiv="refresh"  content="1; url=` + url + `">
	</head>`
	ctx.Render(http.StatusOK, render.Data{
		ContentType: "text/html",
		Data:        []byte(html),
	})
}

// requestForm('post','http://m:4500/redirect/form?redirect_uri=https://s/dump/ab/c/c2',{id_token:"xx",timeout:11,redirect_uri:"https://s/dump/ab/c/c"})
func redirectForm(ctx *gin.Context) {
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

	ctx.HTML(http.StatusOK, "redirect.tmpl", rdata)

}
