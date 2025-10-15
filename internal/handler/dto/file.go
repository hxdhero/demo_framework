package dto

import "lls_api/pkg/util/files"

type RespIdCardDetect struct {
	Results struct {
		BaiduFrontResult map[string]any `json:"baidu_front_result"`
		BaiduBackResult  map[string]any `json:"baidu_back_result"`
	} `json:"results"`
	ServerlessResults string `json:"serverless_results"`
}

func IdCardFrontToResp(front files.IdCardFront, url string) map[string]any {
	res := make(map[string]any)
	res["direction"] = front.Direction
	res["risk_type"] = front.RiskType
	res["idcard_front_image_url"] = url
	res["住址"] = front.WordsResult.Address
	res["公民身份号码"] = front.WordsResult.Number
	res["出生"] = front.WordsResult.Birthday
	res["姓名"] = front.WordsResult.Name
	res["性别"] = front.WordsResult.Gender
	res["民族"] = front.WordsResult.Nation
	res["住址"] = front.WordsResult.Address
	return res
}

func IdCardBackToResp(back files.IdCardBack, url string) map[string]any {
	res := make(map[string]any)
	res["direction"] = back.Direction
	res["risk_type"] = back.RiskType
	res["idcard_back_image_url"] = url
	res["失效日期"] = back.WordsResult.ExpiryDate
	res["签发日期"] = back.WordsResult.IssueDate
	res["签发机关"] = back.WordsResult.IssueBody
	return res
}
