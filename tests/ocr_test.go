package tests

import (
	"bytes"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/stretchr/testify/require"
	"io"
	"lls_api/internal/handler/dto"
	"lls_api/pkg/util/files"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOcrIdcard(t *testing.T) {

	// 获取身份证正反面 todo： mock
	idCardFrontReader, err := files.OSSGet("testcases/lls_go/test_ocr_id_card_front.jpeg", oss.Process("image/resize,w_560,h_720"))
	require.Nil(t, err)
	idCardFrontData, err := io.ReadAll(idCardFrontReader)
	require.Nil(t, err)
	idCardBackReader, err := files.OSSGet("testcases/lls_go/test_ocr_id_card_back_过期.jpeg", oss.Process("image/resize,w_560,h_720"))
	require.Nil(t, err)
	idCardBackData, err := io.ReadAll(idCardBackReader)
	require.Nil(t, err)
	idCardBackPermanentReader, err := files.OSSGet("testcases/lls_go/test_ocr_id_card_back_长期.png", oss.Process("image/resize,w_560,h_720"))
	require.Nil(t, err)
	idCardBackPermanentData, err := io.ReadAll(idCardBackPermanentReader)
	require.Nil(t, err)
	idCardBackValidReader, err := files.OSSGet("testcases/lls_go/test_ocr_id_card_back_有效期内.png", oss.Process("image/resize,w_560,h_720"))
	require.Nil(t, err)
	idCardBackValidData, err := io.ReadAll(idCardBackValidReader)
	require.Nil(t, err)

	// 没有提供身份证反面
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	idcardFront, _ := writer.CreateFormFile("idcard_front", "idcard_front.jpg")
	_, err = io.Copy(idcardFront, bytes.NewReader(idCardFrontData))
	require.Nil(t, err)
	require.Nil(t, writer.Close())
	req, _ := http.NewRequest("POST", "/1/g/1.0/ocr/idcard/", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	service.ServeHTTP(w, req)
	require.Equal(t, 400, w.Code)
	require.Equal(t, "{\"errorMessage\":\"参数异常：idcard_back\"}", w.Body.String())

	// 提供错误文件格式
	b.Reset()
	writer = multipart.NewWriter(&b)
	idcardFront, _ = writer.CreateFormFile("idcard_front", "idcard_front.jpg")
	_, err = io.Copy(idcardFront, bytes.NewReader(idCardFrontData))
	require.Nil(t, err)
	idcardBack, _ := writer.CreateFormFile("idcard_back", "idcard_back.jpg")
	_, err = io.Copy(idcardBack, bytes.NewReader([]byte("test image content")))
	require.Nil(t, err)
	require.Nil(t, writer.Close())
	req, _ = http.NewRequest("POST", "/1/g/1.0/ocr/idcard/", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w = httptest.NewRecorder()
	service.ServeHTTP(w, req)
	require.Equal(t, 400, w.Code)
	require.Equal(t, "{\"errorMessage\":\"获取存储的身份证反面错误\"}", w.Body.String())

	// 身份证过期
	b.Reset()
	writer = multipart.NewWriter(&b)
	idcardFront, _ = writer.CreateFormFile("idcard_front", "idcard_front.jpg")
	_, err = io.Copy(idcardFront, bytes.NewReader(idCardFrontData))
	require.Nil(t, err)
	idcardBack, _ = writer.CreateFormFile("idcard_back", "idcard_back.jpg")
	_, err = io.Copy(idcardBack, bytes.NewReader(idCardBackData))
	require.Nil(t, err)
	require.Nil(t, writer.Close())
	req, _ = http.NewRequest("POST", "/1/g/1.0/ocr/idcard/", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w = httptest.NewRecorder()
	service.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	// 长期身份证
	b.Reset()
	writer = multipart.NewWriter(&b)
	idcardFront, _ = writer.CreateFormFile("idcard_front", "idcard_front.jpg")
	_, err = io.Copy(idcardFront, bytes.NewReader(idCardFrontData))
	require.Nil(t, err)
	idcardBack, _ = writer.CreateFormFile("idcard_back", "idcard_back.jpg")
	_, err = io.Copy(idcardBack, bytes.NewReader(idCardBackPermanentData))
	require.Nil(t, err)
	require.Nil(t, writer.Close())
	req, _ = http.NewRequest("POST", "/1/g/1.0/ocr/idcard/", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w = httptest.NewRecorder()
	service.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	var res dto.RespIdCardDetect
	require.Nil(t, json.Unmarshal([]byte(w.Body.String()), &res))
	require.Equal(t, "长期", res.Results.BaiduBackResult["失效日期"].(map[string]interface{})["words"].(string))

	// 有效期内
	b.Reset()
	writer = multipart.NewWriter(&b)
	idcardFront, _ = writer.CreateFormFile("idcard_front", "idcard_front.jpg")
	_, err = io.Copy(idcardFront, bytes.NewReader(idCardFrontData))
	require.Nil(t, err)
	idcardBack, _ = writer.CreateFormFile("idcard_back", "idcard_back.jpg")
	_, err = io.Copy(idcardBack, bytes.NewReader(idCardBackValidData))
	require.Nil(t, err)
	require.Nil(t, writer.Close())
	req, _ = http.NewRequest("POST", "/1/g/1.0/ocr/idcard/", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w = httptest.NewRecorder()
	service.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)
	require.Nil(t, json.Unmarshal([]byte(w.Body.String()), &res))
	require.Equal(t, "20400729", res.Results.BaiduBackResult["失效日期"].(map[string]interface{})["words"].(string))

}
