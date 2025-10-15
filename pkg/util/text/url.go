package text

import "net/url"

/**
url相关操作
*/

// SafeUnquote URL 编码的字符串转换回原始字符串
func SafeUnquote(s string) string {
	decoded, err := url.QueryUnescape(s)
	if err != nil {
		return s // 出错时返回原始字符串
	}
	return decoded
}
