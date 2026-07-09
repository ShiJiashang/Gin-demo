package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var pageFS embed.FS

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20

	r.GET("/", func(c *gin.Context) {
		html, err := pageFS.ReadFile("index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "read page failed")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", html)
	})

	r.POST("/submit", func(c *gin.Context) {
		dumpPath, err := saveRawHTTPRequest(c)
		if err != nil {
			c.String(http.StatusInternalServerError, "保存原始请求失败: %v", err)
			return
		}

		username := c.PostForm("username")
		age := c.PostForm("age")
		file, _ := c.FormFile("avatar")

		filename := "没有选择文件"
		if file != nil {
			filename = file.Filename
		}

		c.String(http.StatusOK, "提交成功。\n\n原始请求已写入文件:\n%s\n\nusername=%s\nage=%s\navatar=%s\n", dumpPath, username, age, filename)
	})

	fmt.Println("打开浏览器访问: http://localhost:8090")
	if err := r.Run(":8090"); err != nil {
		panic(err)
	}
}

func saveRawHTTPRequest(c *gin.Context) (string, error) {
	req := c.Request

	dump, err := httputil.DumpRequest(req, false)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))

	content := make([]byte, 0, len(dump)+len(body))
	content = append(content, dump...)
	content = append(content, body...)

	dumpDir := filepath.Join("formdemo", "request_dumps")
	if err := os.MkdirAll(dumpDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("request-%s.http", time.Now().Format("20060102-150405.000000000"))
	dumpPath := filepath.Join(dumpDir, filename)
	if err := os.WriteFile(dumpPath, content, 0644); err != nil {
		return "", err
	}

	lastPath := filepath.Join(dumpDir, "last_request.http")
	if err := os.WriteFile(lastPath, content, 0644); err != nil {
		return "", err
	}

	fmt.Println("原始 HTTP 请求已写入:", dumpPath)
	fmt.Println("最新请求副本:", lastPath)

	return dumpPath, nil
}
