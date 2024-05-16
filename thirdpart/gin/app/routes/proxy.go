package routes

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Proxy(ctx *gin.Context) {
	realPath := ctx.Param("path")

	req, err := http.NewRequest(ctx.Request.Method, fmt.Sprintf("http://localhost:8080%v", realPath), ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	for k, vs := range ctx.Request.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("X-Client-Ip", ctx.ClientIP())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(resp.StatusCode)

	header := ctx.Writer.Header()
	for k, vs := range resp.Header {
		for _, v := range vs {
			header.Add(k, v)
		}
	}

	buf := make([]byte, 100*1024)
	for {
		n, err := resp.Body.Read(buf)
		if err == nil || err == io.EOF {
			if _, e := ctx.Writer.Write(buf[:n]); e != nil {
				ctx.AbortWithError(http.StatusInternalServerError, e)
			}
			if err == io.EOF {
				break
			}
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	ctx.Writer.Flush()
}
