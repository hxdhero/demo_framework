package dto

import "github.com/golang-jwt/jwt/v5"

type ReqBlueUserLogin struct {
	AgreeProtocol bool   `json:"agree_protocol"`
	LoginType     int    `json:"login_type"`
	Phone         string `json:"phone"`
	QsApp         int    `json:"qsApp"`
	RelatedId     int    `json:"related_id"`
	SmsCd         string `json:"smsCd"`
}

type RespBlueUserLogin struct {
	ExpAt      jwt.NumericDate `json:"expAt"`
	SaasID     int             `json:"saas_id"`
	Tk         string          `json:"tk"`
	Uid        int             `json:"uId"`
	UserSaasID int             `json:"usersaas_id"`
}
