package files

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"lls_api/pkg/app"
	"lls_api/pkg/config"
	"lls_api/pkg/rdb"
	"lls_api/pkg/rerr"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BaiduAuth struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type BaiduAuthErr struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func refreshAuth() (*BaiduAuth, error) {
	Url := fmt.Sprintf(
		"https://aip.baidubce.com/oauth/2.0/token?client_id=%v&client_secret=%v&grant_type=client_credentials",
		config.C.BaiDuOCR.ClientId, config.C.BaiDuOCR.Secret,
	)
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest("POST", Url, payload)
	if err != nil {
		return nil, rerr.Wrap(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, rerr.Wrap(err)
	}
	defer CloseCloser(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, rerr.Wrap(err)
	}

	var authObj BaiduAuth
	if err = json.Unmarshal(body, &authObj); err != nil {
		return nil, rerr.Wrap(err)
	}
	var authErr BaiduAuthErr
	if authObj.AccessToken == "" {
		if err = json.Unmarshal(body, &authErr); err != nil {
			return nil, rerr.Wrap(err)
		}
		if authErr.Error != "" || authErr.ErrorDescription != "" {
			return nil, rerr.New(fmt.Sprintf("%s: %s", authErr, authErr.ErrorDescription))
		}
		return nil, rerr.New("获取baidu ocr token 失败.")
	}

	return &authObj, nil
}

func doRequest(url string, payload url.Values, auth *BaiduAuth) ([]byte, error) {
	newUrl := fmt.Sprintf("%v?access_token=%v", url, auth.AccessToken)
	client := &http.Client{}

	req, err := http.NewRequest("POST", newUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, rerr.Wrap(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, rerr.Wrap(err)
	}
	defer CloseCloser(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, rerr.Wrap(err)
	}

	if resp.StatusCode != 200 {
		return nil, rerr.New(string(data))
	}

	return data, nil
}

// IdCardFile 身份证文件
type IdCardFile struct {
	Front []byte // 正面文件
	Back  []byte // 反面文件
}

type IdCardFieldLocation struct {
	Height int `json:"height"`
	Left   int `json:"left"`
	Top    int `json:"top"`
	Width  int `json:"width"`
}

type IdCardFront struct {
	IdCardNumberType int    `json:"idcard_number_type"`
	ImageStatus      string `json:"image_status"`
	WordsResult      struct {
		Address struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"住址"`
		Number struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"公民身份号码"`
		Birthday struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"出生"`
		Name struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"姓名"`
		Gender struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"性别"`
		Nation struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"民族"`
	} `json:"words_result"`
	Direction int    `json:"direction"` // 图像方向
	RiskType  string `json:"risk_type"` // 身份证风险类型:normal-正常身份证；copy-复印件；temporary-临时身份证；screen-翻拍；screenshot-截屏（仅在开启 detect_screenshot 时返回）；unknown-其他未知情况
}

type IdCardBack struct {
	ImageStatus string `json:"image_status"`
	WordsResult struct {
		ExpiryDate struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"失效日期"`
		IssueDate struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"签发日期"`
		IssueBody struct {
			Location IdCardFieldLocation `json:"location"`
			Words    string              `json:"words"`
		} `json:"签发机关"`
	} `json:"words_result"`
	Direction int    `json:"direction"` // 图像方向
	RiskType  string `json:"risk_type"` // 身份证风险类型:normal-正常身份证；copy-复印件；temporary-临时身份证；screen-翻拍；screenshot-截屏（仅在开启 detect_screenshot 时返回）；unknown-其他未知情况
}

// IdCardOCRResult 身份证识别结果
type IdCardOCRResult struct {
	Front IdCardFront `json:"Front"`
	Back  IdCardBack  `json:"Back"`
}

type OCRError struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// OcrIdCardDiscern ocr身份证识别
func OcrIdCardDiscern(ctx *app.Context, idCard IdCardFile, auth *BaiduAuth, isRaise bool) (IdCardOCRResult, error) {
	var result IdCardOCRResult
	baseUrl := "https://aip.baidubce.com/rest/2.0/ocr/v1/idcard"
	payload := make(url.Values)
	payload.Set("detect_direction", "true")
	payload.Set("detect_risk", "true")

	// 正面识别
	payload.Set("image", base64.StdEncoding.EncodeToString(idCard.Front))
	payload.Set("id_card_side", "front")
	frontResp, err := doRequest(baseUrl, payload, auth)
	if err != nil {
		return result, rerr.WrapS(err, "身份证人像面识别失败")
	}
	var frontResult IdCardFront
	if err := json.Unmarshal(frontResp, &frontResult); err != nil {
		return result, rerr.WrapS(err, "身份证人像面识别失败")
	}
	if frontResult.WordsResult.Name.Words == "" {
		var ocrErr OCRError
		if err = json.Unmarshal(frontResp, &ocrErr); err != nil {
			return result, rerr.WrapS(err, "识别身份证人像面失败")
		}
		return result, rerr.New("识别身份证人像面失败")
	}

	// 反面识别
	payload.Set("image", base64.StdEncoding.EncodeToString(idCard.Back))
	payload.Set("id_card_side", "back")
	backResp, err := doRequest(baseUrl, payload, auth)
	if err != nil {
		return result, rerr.WrapS(err, "身份证国徽面识别失败")
	}
	var backResult IdCardBack
	if err := json.Unmarshal(backResp, &backResult); err != nil {
		return result, rerr.WrapS(err, "身份证国徽面识别失败")
	}
	if backResult.WordsResult.IssueDate.Words == "" {
		var ocrErr OCRError
		if err = json.Unmarshal(backResp, &ocrErr); err != nil {
			return result, rerr.WrapS(err, "身份证国徽面识别失败")
		}
		return result, rerr.New("身份证国徽面识别失败")
	}
	if isRaise {
		if backResult.WordsResult.ExpiryDate.Words != "长期" {
			t, err := time.Parse("20060102", backResult.WordsResult.ExpiryDate.Words)
			if err != nil {
				ctx.Log().Errorf("无法识别身份证的有效期限:[%s]", backResult.WordsResult.ExpiryDate.Words)
				return result, rerr.New("无法识别身份证的有效期限")
			}
			if t.Before(time.Now()) {
				ctx.Log().Errorf("身份证的有效期限已过:[%s]", backResult.WordsResult.ExpiryDate.Words)
				return result, rerr.New("身份证的有效期限已过")
			}
		}
	}

	result.Front = frontResult
	result.Back = backResult
	return result, nil
}

// OcrAuth 从缓存获取ocrAuth
func OcrAuth(ctx *app.Context, redisKey string) (*BaiduAuth, error) {
	// redis 获取缓存的token
	bs, err := ctx.Redis().Get(ctx.Context(), redisKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, rerr.Wrap(err)
	}

	var auth *BaiduAuth
	if err == nil {
		// Redis 中有数据
		auth = &BaiduAuth{}
		if err := json.Unmarshal(bs, auth); err != nil {
			return nil, rerr.Wrap(err)
		}
	} else {
		// Redis 中没有数据，需要重新获取
		ctx.Log().Infof("redis key:%s 不存在,发起请求获取token", rdb.FullKey(redisKey))
		auth, err = refreshAuth()
		if err != nil {
			ctx.Log().Errorf("请求baidu token错误:%s", err.Error())
			return nil, rerr.Wrap(err)
		}

		// 缓存新获取的token
		bs, err = json.Marshal(auth)
		if err != nil {
			return nil, rerr.Wrap(err)
		}
		// 设置过期时间为提前一天
		expire := time.Duration(auth.ExpiresIn)*time.Second - time.Hour*24
		if err := ctx.Redis().Set(ctx.Context(), redisKey, bs, expire).Err(); err != nil {
			return nil, rerr.Wrap(err)
		}
		ctx.Log().Infof("获取到新token,过期时间 %v", auth.ExpiresIn)
	}

	return auth, nil
}
