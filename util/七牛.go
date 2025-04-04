package util

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/tidwall/gjson"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

// Upload 上传到七牛
func F上传(localFile string, keyToOverwrite string, expires uint64) (string, error) {
	qiniu := GetQiniu("iryna")
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", qiniu.Bucket, keyToOverwrite),
	}
	if expires == 0 {
		expires = 3600 * 24 * 365 * 20 //20年
	}
	putPolicy.Expires = expires //1小时有效期
	mac := auth.New(qiniu.AccessKey, qiniu.SecrectKey)
	upToken := putPolicy.UploadToken(mac)

	formUploader := storage.NewFormUploader(&qiniu.Config)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, keyToOverwrite, localFile, &putExtra)
	if err != nil {
		return "", err
	}
	return qiniu.Domain + keyToOverwrite, nil
}

type C七牛结构 struct {
	Bucket     string
	AccessKey  string
	SecrectKey string
	Domain     string
	Config     storage.Config
}

// GetQiniu 链接七牛
func GetQiniu(drive string) C七牛结构 {
	buf, _ := ioutil.ReadFile("db.json")
	s := gjson.Get(string(buf), drive)
	qiniu := C七牛结构{}
	qiniu.Bucket = s.Get("bucket").String()
	qiniu.AccessKey = s.Get("accesskey").String()
	qiniu.SecrectKey = s.Get("secretkey").String()
	qiniu.Domain = s.Get("domain").String()
	qiniu.Config = storage.Config{
		UseHTTPS:      true,
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
	}
	return qiniu
}
