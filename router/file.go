package router

import (
	"fmt"
	"ginapp/nettool"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const UPLOAD_DIR = "./tmp"

// curl http://m:4500/f/r/my/a.log -o a.log
func fileReadHandler(c *gin.Context) {
	filepath := strings.TrimPrefix(c.Param("path"), "/")
	if isBadPath(filepath) {
		panic("bad file path: " + filepath)
	}
	if strings.HasSuffix(filepath, ".json") {
		c.Writer.Header().Set("Content-type", "application/json")
	} else if strings.Contains(filepath, ".txt") {
		c.Writer.Header().Set("Content-type", "text/plain; charset=utf-8")
	} else if strings.Contains(filepath, ".mp4") {
		c.Writer.Header().Set("Content-type", "video/mp4")
	} else if strings.Contains(filepath, ".html") {
		c.Writer.Header().Set("Content-type", "text/html; charset=utf-8")
	} else {
		c.Writer.Header().Set("Content-type", "application/octet-stream")
	}
	f, err := os.Open(UPLOAD_DIR + "/" + filepath)
	if err != nil {
		panic(err)
	}
	io.Copy(c.Writer, f)
}

// curl m:4500/f/w/my/a.log -F file1=@test.log -F name=alex | jq
func fileWriteHandler(c *gin.Context) {
	name := c.PostForm("name")
	upload_path := strings.TrimPrefix(c.Param("path"), "/")
	file, err := c.FormFile("file1")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	} else {
		if upload_path == "" {
			upload_path = file.Filename
		}
		if isBadPath(upload_path) {
			c.String(http.StatusBadRequest, fmt.Sprintf("bad upload path: %s", upload_path))
			return
		}
		// filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, UPLOAD_DIR+"/"+upload_path); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["files"]
	filelist := []string{}
	for _, file := range files {
		filelist = append(filelist, file.Filename)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     gin.H{"name": name},
		"file1":    filepath.Base(file.Filename),
		"filelist": filelist,
		"msg":      fmt.Sprintf("download: curl %s:4500/f/r/%s  -o a.log", nettool.GetLocalIP(), upload_path),
	})
}

var bad_path_pattern = regexp.MustCompile(`\./|//|\.\.`)
var valid_path_pattern = regexp.MustCompile(`^[\w\-\./]+$`)

func isBadPath(path string) bool {
	if bad_path_pattern.MatchString(path) {
		return true
	}
	if !valid_path_pattern.MatchString(path) {
		return true
	}
	if strings.HasPrefix(path, "/") {
		return true
	}
	return false

}
