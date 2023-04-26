package common

import (
	"io"
	"nftserver/global"
	"nftserver/pkg/app"
	"nftserver/pkg/utils"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type CommonAPI struct{}

func NewCommonAPI() CommonAPI {
	return CommonAPI{}
}

// @Router /v1/common/images/upload [post]
func (u CommonAPI) ImagesUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(app.RespErr(app.ErrForm, "File not found"))
		return
	}

	conf := global.StorageSetting

	var maxFileSize = conf.Images.MaxSize << 20
	if file.Size > int64(maxFileSize) {
		c.JSON(app.RespErr(app.ErrFileSize, ""))
		return
	}

	isAllow := false
	extstring := strings.ToUpper(path.Ext(file.Filename))
	for _, v := range conf.Images.Types {
		if extstring == v {
			isAllow = true
			break
		}
	}
	if !isAllow {
		c.JSON(app.RespErr(app.ErrFileType, ""))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(app.RespErr(app.ErrFileFile, err.Error()))
		return
	}
	defer src.Close()

	fileHash := utils.SHA1FileHash(file)
	newFileName := fileHash + extstring
	storagePath := path.Join(conf.BasePath, conf.Images.StorePath, newFileName)

	out, err := os.Create(storagePath)
	if err != nil {
		c.JSON(app.RespErr(app.ErrFileFile, err.Error()))
		return
	}
	defer out.Close()
	io.Copy(out, src)

	newFileName = path.Join("/", conf.Images.StorePath, newFileName)
	c.JSON(app.RespOK(newFileName))
}

// @Router /v1/common/resources/upload [post]
func (u CommonAPI) ResourcesUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(app.RespErr(app.ErrForm, "File not found"))
		return
	}

	conf := global.StorageSetting

	var maxFileSize = conf.Resources.MaxSize << 20
	if file.Size > int64(maxFileSize) {
		c.JSON(app.RespErr(app.ErrFileSize, ""))
		return
	}

	isAllow := false
	extstring := strings.ToUpper(path.Ext(file.Filename))
	for _, v := range conf.Resources.Types {
		if extstring == v {
			isAllow = true
			break
		}
	}
	if !isAllow {
		c.JSON(app.RespErr(app.ErrFileType, ""))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(app.RespErr(app.ErrFileFile, err.Error()))
		return
	}
	defer src.Close()

	fileHash := utils.SHA1FileHash(file)
	newFileName := fileHash + extstring
	storagePath := path.Join(conf.BasePath, conf.Resources.StorePath, newFileName)

	out, err := os.Create(storagePath)
	if err != nil {
		c.JSON(app.RespErr(app.ErrFileFile, err.Error()))
		return
	}
	defer out.Close()
	io.Copy(out, src)

	newFileName = path.Join("/", conf.Resources.StorePath, newFileName)
	c.JSON(app.RespOK(newFileName))
}
