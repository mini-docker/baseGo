package controllers

import (
	"baseGo/src/fecho/echo"
	"baseGo/src/fecho/golog"
	"baseGo/src/model/code"
	"baseGo/src/red_agency/app/controller/common"
	"baseGo/src/red_agency/app/middleware/validate"
	"baseGo/src/red_agency/app/server"
	"baseGo/src/red_agency/app/services/pubresource"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	dividingLine   = "/"
	fileKeyNotFind = errors.New(`not find 'file' key in form-data`)
)

type UploadFileController struct {
}

type FileData struct {
	Content       []byte              `json:"content"`
	ContentType   string              `json:"content_type"`
	ContentLength int64               `json:"content_length"`
	MetaData      map[string][]string `json:"meta_data"`
}

type JsonReturn struct {
	Msg  string   `json:"msg"`
	Data FileData `json:"data"`
}

// 文件上传
func (ac UploadFileController) UpLoadFile(ctx server.Context) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		golog.Error("AccountController", "GoIMUpdate", "error:", err)
		return common.HttpResultJsonError(ctx, &validate.Err{Code: code.FILE_CAN_NOT_BE_EMPTY})
	}
	files := form.File["file"]
	if len(files) == 0 {
		return common.HttpResultJsonError(ctx, fileKeyNotFind)
	}

	fileHandle, err := files[0].Open()
	if err != nil {
		golog.Error("AccountController", "GoIMUpdate", "error:", err)
		return common.HttpResultJsonError(ctx, err)
	}
	defer fileHandle.Close()

	fileContentByUpload, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		golog.Error("AccountController", "GoIMUpdate", "error:", err)
		return common.HttpResultJsonError(ctx, err)
	}

	filePath, respHeader, err := common.UploadFile(files[0].Filename, fileContentByUpload)
	if err != nil {
		return common.HttpResultJsonError(ctx, err)
	}

	var thumbPath string
	ext := filepath.Ext(filePath)
	ext = strings.ToLower(ext)

	if ext == ".mp4" || ext == ".mov" {
		thumbPath = filePath[:len(filePath)-len(ext)] + "_thumb.jpg"
	} else {
		if ext != ".wav" {
			thumbPath = filePath[:len(filePath)-len(ext)] + "_thumb" + ext
		} else {
			thumbPath = ""
		}

	}
	resMap := make(map[string]interface{}, 0)
	resMap["ext"] = ext
	resMap["duration"] = respHeader.Get("Duration")
	resMap["size"] = respHeader.Get("Size")
	resMap["thumbWidth"] = respHeader.Get("Thumbwidth")
	resMap["thumbHeight"] = respHeader.Get("Thumbheight")
	resMap["width"] = respHeader.Get("Width")
	resMap["height"] = respHeader.Get("Height")
	resMap["url"] = filePath
	resMap["thumbUrl"] = thumbPath
	return common.HttpResultJson(ctx, resMap)
}

// 文件下载
func (ac UploadFileController) DownLoadFile(ctx server.Context) error {
	// 读取url文件路径
	fileFullPath := ctx.Param("*")

	// 处理路径
	fileFullPath = fmt.Sprintf("%s%s", dividingLine, filepath.Clean(fileFullPath))

	// 下载文件
	resp, err := pubresource.NewPublicResourceService().Bucket(ioutil.NopCloser(strings.NewReader(fmt.Sprintf("{\"action\":\"getContent\",\"path\":\"%v\"}", fileFullPath))))
	if err != nil {
		fmt.Println("下载失败:", err)
		golog.Error("UploadFileController", "DownLoadFile", "download error:", err)
		return ctx.HTML(404, "")
	}
	// 读取文件二进制
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取文件二进制失败:", err)
		golog.Error("UploadFileController", "DownLoadFile", "download error:", err)
		return ctx.HTML(404, "")
	}
	jsonData := new(JsonReturn)
	json.Unmarshal(data, &jsonData)
	if nil == jsonData {
		golog.Error("UploadFileController", "DownLoadFile", "download error:json is nil", nil)
		return ctx.HTML(404, "")
	}

	ctx.Response().Header().Set(echo.HeaderContentLength, fmt.Sprint(jsonData.Data.ContentLength))

	ctx.Response().Header().Set("Cache-Control", "max-age=29030400,public")
	return ctx.Blob(http.StatusOK, jsonData.Data.ContentType, jsonData.Data.Content)
}
