package router

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"ginapp/test"

	"github.com/ahuigo/gohttptool/httpreq"
)

func TestAddUser(t *testing.T) {
	// build request
	data := User{
		Name: "testuser",
		City: "Beijing",
	}
	req, _ := httpreq.R().SetReq("POST","/api/v1/user").SetJson(data).ToRequest()

	// send request
	respRecorder, ctx := test.CreateTestCtx(req)
	AddUser(ctx)

	// test response status
	if respRecorder.Code != http.StatusOK {
		errors := ctx.Errors.Errors()
		fmt.Println("output", errors)
		t.Fatal("Expect code 200, but get ", respRecorder.Code, "reason", respRecorder.Body)
	}

	// test response body
	recvUser := User{}
	resp := httpreq.BuildResponse(respRecorder.Result())
	resp.Json(&recvUser)
	if recvUser.Name != data.Name || recvUser.City != data.City {
		t.Fatalf("unexpected response:%v", resp.Text())
	}

}

func TestConfApi(t *testing.T) {
	req, _ := httpreq.R().ToRequest()

	// send request
	respRecorder, ctx := test.CreateTestCtx(req)
	confApi(ctx)

	// test response status
	if respRecorder.Code != http.StatusOK {
		errors := ctx.Errors.Errors()
		fmt.Println("output", errors)
		t.Fatal("Expect code 200, but get ", respRecorder.Code, "reason", respRecorder.Body)
	}

	// test response body
	resp := httpreq.BuildResponse(respRecorder.Result())
	text := resp.Text()
	if !strings.Contains(text, `Port":543`) {
		t.Fatalf("unexpected response:%v", text)
	}
	t.Log(text)

}
