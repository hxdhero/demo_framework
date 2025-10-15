package text

import (
	"lls_api/pkg/rerr"
	"strings"
)

/**
文本提取相关
*/

// ExtraAgeFromIdCardNO 从身份证号中提取出身年月日字符串
func ExtraAgeFromIdCardNO(IdCardNO string) (string, error) {
	IdCardNO = strings.TrimSpace(IdCardNO)

	var birthStr string
	switch len(IdCardNO) {
	case 18:
		birthStr = IdCardNO[6:14] // YYYYMMDD
	case 15:
		birthStr = "19" + IdCardNO[6:12] // 补全年份为19YYMMDD
	default:
		return "", rerr.New("无效的身份证长度")
	}
	return birthStr, nil
}
