package common

import "github.com/golang-jwt/jwt/v5"

type Auth struct {
	Exp        jwt.NumericDate `json:"exp"`
	InstanceId int             `json:"instance_id"`
	QsApp      int             `json:"qsApp"`
	DrfUser    string          `json:"drf_user"`
	Uid        int             `json:"uId"`
	UserId     int             `json:"userId"`
	UserType   string          `json:"userType"`
	IsSuper    bool            `json:"is_super"`
	Extra      map[string]any  `json:"extra"`
}

type UserInstance interface {
	EmployerField() string
	GetStatus() UserAgencyStatus
}

type UserAgencyStatus int

const (
	UserAgencyStatusEnabled  UserAgencyStatus = 1
	UserAgencyStatusDisabled UserAgencyStatus = 2
)

func (s UserAgencyStatus) String() string {
	switch s {
	case UserAgencyStatusEnabled:
		return "启用中"
	case UserAgencyStatusDisabled:
		return "已禁用"
	default:
		return "未知"
	}
}
