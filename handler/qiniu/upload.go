package qiniu

import (
	"fmt"
	"net/http"

	"github.com/1024casts/1024casts/pkg/app"
	"github.com/1024casts/1024casts/pkg/errno"
	"github.com/1024casts/1024casts/service"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

func UploadImage(c *gin.Context) {

	// single file
	file, err := c.FormFile("file")
	if err != nil {
		log.Warnf("[upload] get file err: %v", err)
		app.Response(c, errno.ErrGetUploadFile, nil)
		return
	}

	qiNiuSrv := service.NewQiNiuService()
	uploadRet, err := qiNiuSrv.UploadImage(c, file, true)
	if err != nil {
		log.Warnf("[upload] upload file err: %v", err)
		app.Response(c, errno.ErrUploadingFile, nil)
		return
	}
	fmt.Println(uploadRet)

	app.Response(c, nil, uploadRet)

}

func UploadVideo(c *gin.Context) {

	// single file
	file, err := c.FormFile("file")
	if err != nil {
		log.Warnf("[upload] get file err: %v", err)
		app.Response(c, errno.ErrGetUploadFile, nil)
		return
	}

	qiNiuSrv := service.NewQiNiuService()
	uploadRet, err := qiNiuSrv.UploadVideo(c, file, false)
	if err != nil {
		log.Warnf("[upload] upload file err: %v", err)
		app.Response(c, errno.ErrUploadingFile, nil)
		return
	}
	fmt.Println(uploadRet)

	app.Response(c, nil, uploadRet)

}

// 论坛编辑器上传调用，没有token加密
func WebUpload(c *gin.Context) {

	// single file
	file, err := c.FormFile("file")
	if err != nil {
		log.Warnf("[upload] get file err: %v", err)
		app.Response(c, errno.ErrGetUploadFile, nil)
		return
	}

	qiNiuSrv := service.NewQiNiuService()
	uploadRet, err := qiNiuSrv.UploadImage(c, file, true)
	if err != nil {
		log.Warnf("[upload] upload file err: %v", err)
		app.Response(c, errno.ErrUploadingFile, nil)
		return
	}
	fmt.Println(uploadRet)

	c.JSON(http.StatusOK, gin.H{
		"filename": uploadRet.Url,
	})

}
