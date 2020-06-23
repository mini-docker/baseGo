package pubresource

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/red_agency/conf"
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sync"
)

type publicResourceService struct {
	host      string
	accessKey string
	secretKey string
	c         *http.Client
}

var defaultPublicResourceServiceInstance *publicResourceService

var onceByNewPublicResourceService sync.Once

func NewPublicResourceService() *publicResourceService {
	onceByNewPublicResourceService.Do(func() {
		defaultPublicResourceServiceInstance = new(publicResourceService)
		defaultPublicResourceServiceInstance.host = conf.GetStorageServConfig().Host
		defaultPublicResourceServiceInstance.accessKey = conf.GetStorageServConfig().AccessKey
		defaultPublicResourceServiceInstance.secretKey = conf.GetStorageServConfig().SecretKey
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		defaultPublicResourceServiceInstance.c = http.DefaultClient
	})
	return defaultPublicResourceServiceInstance
}

func (prs *publicResourceService) getHttpRequestByJson(reqData io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, prs.host+"/api/storage", reqData)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("SecretKey", prs.secretKey)
	req.Header.Set("AccessKey", prs.accessKey)
	return req, nil
}

func (prs *publicResourceService) getHttpRequestByFile(reqData *http.Request) (*http.Request, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	err := reqData.ParseMultipartForm(1024 * 1024 * 4)
	if err != nil {
		return nil, err
	}

	// add filed 'dest'
	fwByDestField, _ := w.CreateFormField("dest")
	fwByDestField.Write([]byte(reqData.MultipartForm.Value["dest"][0]))

	// add filed 's'
	fwBySiteIdField, _ := w.CreateFormField("s")
	fwBySiteIdField.Write([]byte(reqData.MultipartForm.Value["s"][0]))

	// add filed 'i'
	fwBySiteIndexIdField, _ := w.CreateFormField("i")
	fwBySiteIndexIdField.Write([]byte(reqData.MultipartForm.Value["i"][0]))

	// add filed 'a'
	fwByAbsolutePathField, _ := w.CreateFormField("a")
	fwByAbsolutePathField.Write([]byte("1")) // ignore 's' and 'i' filed

	// add filed 'file'
	fw, err := w.CreateFormFile("file", reqData.MultipartForm.File["file"][0].Filename)
	if err != nil {
		return nil, err
	}
	remoteFile, err := reqData.MultipartForm.File["file"][0].Open()
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, remoteFile)
	if err != nil {
		return nil, err
	}
	remoteFile.Close()
	w.Close()

	req, err := http.NewRequest(http.MethodPost, prs.host+"/api/storage/upload", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("SecretKey", prs.secretKey)
	req.Header.Set("AccessKey", prs.accessKey)
	return req, nil
}

func (prs *publicResourceService) getHttpRequestByZip(reqData *http.Request) (*http.Request, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	err := reqData.ParseMultipartForm(1024 * 1024 * 4)
	if err != nil {
		return nil, err
	}

	// add filed 'dest'
	fwByDestField, _ := w.CreateFormField("dest")
	fwByDestField.Write([]byte(reqData.MultipartForm.Value["dest"][0]))

	remoteFile, err := reqData.MultipartForm.File["file"][0].Open()
	if err != nil {
		return nil, err
	}

	// add filed 'etag'
	newDest := bytes.NewBuffer(nil)
	_, err = io.Copy(newDest, remoteFile)
	if err != nil {
		return nil, err
	}
	fwByAbsolutePathField, _ := w.CreateFormField("etag")
	fwByAbsolutePathField.Write([]byte(fmt.Sprintf("%x", md5.Sum(newDest.Bytes()))))

	// add filed 'file'
	fw, err := w.CreateFormFile("file", reqData.MultipartForm.File["file"][0].Filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, bytes.NewReader(newDest.Bytes()))
	if err != nil {
		return nil, err
	}
	remoteFile.Close()
	w.Close()

	req, err := http.NewRequest(http.MethodPost, prs.host+"/api/storage/uploadTpl", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("SecretKey", prs.secretKey)
	req.Header.Set("AccessKey", prs.accessKey)
	return req, nil
}

func (prs *publicResourceService) Bucket(reqData io.ReadCloser) (*http.Response, error) {
	defer reqData.Close()

	reqBody, err := ioutil.ReadAll(reqData)
	if err != nil {
		golog.Error("publicResourceService", "Bucket", "err:", err)
		return nil, err
	}

	type actionCheck struct {
		Action string `json:"action"`
	}
	var ac actionCheck
	err = json.Unmarshal(reqBody, &ac)
	if err != nil {
		golog.Error("publicResourceService", "Bucket", "err:", err)
		return nil, err
	}
	switch ac.Action {
	case "removeBucket", "hardRemove":
		golog.Error("publicResourceService", "Bucket", "err:", errors.New(`action 'removeBucket' and 'hardRemove' is not allow`))
		return nil, err
	}

	req, err := prs.getHttpRequestByJson(bytes.NewReader(reqBody))
	if err != nil {
		golog.Error("publicResourceService", "Bucket", "err:", err)
		return nil, err
	}
	return prs.c.Do(req)
}

func (prs *publicResourceService) BucketUploadByFile(reqData *http.Request) (*http.Response, error) {
	req, err := prs.getHttpRequestByFile(reqData)
	if err != nil {
		golog.Error("publicResourceService", "BucketUploadByFile", "err:", err)
		return nil, err
	}
	return prs.c.Do(req)
}

func (prs *publicResourceService) BucketUploadByZip(reqData *http.Request) (*http.Response, error) {
	req, err := prs.getHttpRequestByZip(reqData)
	if err != nil {
		golog.Error("publicResourceService", "BucketUploadByZip", "err:", err)
		return nil, err
	}
	return prs.c.Do(req)
}
