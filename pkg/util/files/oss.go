package files

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"lls_api/pkg/config"
	"lls_api/pkg/log"
	"lls_api/pkg/rerr"
	"strings"
)

var (
	ossClient *oss.Client
)

func InitOss() {
	if ossClient != nil {
		return
	}
	config.C.OSS.Prefix = fmt.Sprintf("https://%s.%s/", config.C.OSS.Bucket, config.C.OSS.Endpoint)
	log.DefaultContext().Infof("初始化 bucket:%s, prefix:%s ", config.C.OSS.Bucket, config.C.OSS.Prefix)
	var err error
	ossClient, err = oss.New(config.C.OSS.Endpoint, config.C.OSS.AccessKeyID, config.C.OSS.AccessKeySecret)
	if err != nil {
		panic(err)
	}
}

func OSSPut(objectKey string, reader io.Reader, options ...oss.Option) error {
	bucket, err := ossClient.Bucket(config.C.OSS.Bucket)
	if err != nil {
		return rerr.Wrap(err)
	}
	err = bucket.PutObject(objectKey, reader, options...)
	return rerr.Wrap(err)
}

func OSSGet(objectKey string, options ...oss.Option) (io.ReadCloser, error) {
	bucket, err := ossClient.Bucket(config.C.OSS.Bucket)
	if err != nil {
		return nil, rerr.Wrap(err)
	}
	reader, err := bucket.GetObject(objectKey, options...)
	if err != nil {
		return nil, rerr.Wrap(err)
	}
	return reader, nil
}

func OssSignURL(objectKey string, expiredSec int64) (signedURL string, err error) {
	bucket, err := ossClient.Bucket(config.C.OSS.Bucket)
	if err != nil {
		return signedURL, err
	}
	objectKey = strings.Split(objectKey, "?")[0]
	if strings.HasPrefix(objectKey, "http") {
		objectKey = strings.Split(objectKey, ".com/")[1]
	}
	signedURL, err = bucket.SignURL(objectKey, oss.HTTPGet, expiredSec)
	if err != nil {
		return signedURL, err
	}
	return signedURL, nil
}
