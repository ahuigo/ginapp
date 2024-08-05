package middleware

import (
	"ginapp/test"
	"io"
	"testing"

	"github.com/ahuigo/gohttptool/httpreq"
)


func TestGetUsername(t *testing.T) {
	req,_ := httpreq.R().SetJson(map[string]any{
		"username": "alex",
	}).SetReq("POST","http://m/api/v1/mauth/login").ToRequest()
	curl := httpreq.BuildCurlCommand(req, nil)
	t.Log(curl)

	resp, ctx := test.CreateTestCtx(req)
	username1 := getLoginUserFromBody(ctx)
	if username1 != "alex" {
		bytes, _ := io.ReadAll(resp.Result().Body)
		t.Fatalf("UnExpected name:%s, body:%s", username1, bytes)
	}
}