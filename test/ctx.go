package test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func CreateTestCtx(req *http.Request) (resp *httptest.ResponseRecorder, ctx *gin.Context) {
	resp = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(resp)
	ctx.Request = req
	return
}
