package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IndexHandler index handler
func IndexHandler(c *gin.Context) {
	// 获取 GET 请求中的参数信息
	if err := c.Request.ParseForm(); err != nil {
		// 输出 json 结果给调用方
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("parse form failed: %v", err),
		})
		return
	}
	param, err := getParam(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("get param failed: %v", err),
		})
		return
	}


	return
}

// getParam 获取参数信息
func getParam(r *http.Request) (Param, error) {
	if err := r.ParseForm(); err != nil {
		return Param{}, err
	}
	return Param{
		CurrentContext: r.Form.Get("context"),
		Namespace:      r.Form.Get("namespace"),
		Pod:            r.Form.Get("pod"),
		Container:      r.Form.Get("container"),
	}, nil
}
