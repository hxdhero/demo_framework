package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"lls_api/common"
	"lls_api/internal"
	"lls_api/internal/model"
	"lls_api/internal/model/base"
	"lls_api/internal/model/gen"
	"lls_api/pkg"
	"lls_api/pkg/app"
	"lls_api/pkg/config"
	"lls_api/pkg/util/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	ctx             *app.Context
	service         app.Server
	defaultUsersaas model.UserSaas
)

func TestMain(m *testing.M) {
	// 初始化依赖
	pkg.InitGlobal()
	// 初始化appContext
	ctx = app.ContextWithTest()
	// 表结构迁移
	if err := model.MigrateModels(ctx); err != nil {
		panic(err)
	}
	// 初始化server
	services, err := internal.Servers()
	if err != nil {
		panic(err)
	}
	service = services[0]
	app.WithTestMTransaction(m, ctx, func(txCtx *app.Context) error {
		t := time.Now()
		// 初始化 user
		user := model.User{BlueUser: gen.BlueUser{Mobile: "18989898989", CreateAt: t, ModifyAt: t}}
		if err := user.Create(txCtx); err != nil {
			panic(err)
		}
		// 初始化 saas
		saas := model.Saas{BlueSaas: gen.BlueSaas{Name: "测试saas", CreateAt: t, ModifyAt: t}}
		if err := saas.Create(txCtx); err != nil {
			panic(err)
		}
		// 初始化address
		// address := model.Address{Address:gen.Address{Province: "上海市",City: "上海市",District: "闵行区",CreateAt: t,ModifyAt: t}}
		// 初始化company
		// company := model.BlueCompany{BlueCompany:gen.BlueCompany{IsEnable: true,RegisterAddress:"上海市xxx路",SourceType: 1,SettlementDay: 10,AutowithdrawStatus: 1,SettlementType: 1,CreateAt: t,ModifyAt: t}}

		// 初始化 usersaas
		defaultUsersaas = model.UserSaas{BlueUsersaas: gen.BlueUsersaas{UserID: user.ID, SaasID: saas.ID, CreateAt: t, ModifyAt: t}}
		if err := defaultUsersaas.Create(txCtx); err != nil {
			panic(err)
		}
		app.SetTestTransaction(txCtx)
		ctx = txCtx
		m.Run()
		return nil
	})
}

func NewReqWithUser(t *testing.T, user model.UserSaas, method string, url string, body io.Reader) (*http.Request, error) {
	t.Helper()
	auth := common.Auth{
		Exp:        jwt.NewJwtTime(time.Now().Add(time.Duration(config.C.JWT.Exp) * time.Second)),
		InstanceId: int(user.InstanceId()),
		QsApp:      1,
		Uid:        int(user.GetID()),
		UserId:     int(user.UserID),
	}
	token, err := jwt.GenerateToken(auth, config.C.JWT.Secret)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("HTTP_AUTHORIZATION", token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}

func ResDisplayErr(t *testing.T, w *httptest.ResponseRecorder) common.DisplayError {
	var res common.DisplayError
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	return res
}

func setToken(req *http.Request, token string) {
	req.Header.Set("HTTP_AUTHORIZATION", token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

type TestUser struct {
}

func CreateCDLAgency(ctx *app.Context, t *testing.T, name string, superManagerMobile string, loginMobile string, laborCertificationTypes []int, requirementAgencyId int, platformId int) *TestUser {
	// attorney
	attorney := model.Attorney{BlueAttorney: gen.BlueAttorney{Name: "1", Mobile: "13520240125", LetteryURL: "2", IDURL: "3", IDNo: "4"}}
	assert.NoError(t, ctx.DB().Create(&attorney).Error)
	// address
	var address model.Address
	assert.NoError(t, ctx.DB().Last(&address))
	// domain
	var domain model.Domain
	assert.NoError(t, ctx.DB().First(&domain, "platform_id = ?", platformId))
	// agency
	bs, err := json.Marshal(laborCertificationTypes)
	assert.NoError(t, err)
	agency := model.Agency{
		BlueAgency: gen.BlueAgency{
			ZhaoshangSettlementChannel: true,
			LaborCertificationTypes:    string(bs),
			IsCompleted:                true,
			IsEnable:                   true,
			RequirementAgencyID:        base.NullIDFromInt(requirementAgencyId),
			AddressID:                  address.ID,
			AttorneyID:                 attorney.ID,
			Telephone:                  "18888888888",
			PlatformID:                 base.NullIDFromInt(platformId),
			AgencyBusinessModel:        "1000",
			SourceType:                 2,
			TaxCollectType:             1,
			Name:                       name,
			DomainID:                   base.NullIDFromID(domain.ID),
			SuperManagerMobile:         superManagerMobile,
		},
	}
	assert.NoError(t, ctx.DB().Create(&agency).Error)
	// user
	managerUser := model.User{BlueUser: gen.BlueUser{Mobile: superManagerMobile}}
	assert.NoError(t, ctx.DB().Create(&managerUser).Error)
	// user agency
	managerUserAgency := model.UserAgency{BlueUseragency: gen.BlueUseragency{Name: "陈娟3", UserID: base.NullIDFromID(managerUser.ID), AgencyID: agency.ID, IsSuper: true}}
	assert.NoError(t, ctx.DB().Create(&managerUserAgency).Error)
	/*// agencyLimitLaborAge
	assert.NoError(t, ctx.DB().Create(model.AgencyLimitLaborAge{BlueAgencylimitlaborage:gen.BlueAgencylimitlaborage{AgencyID: agency.ID,MaxFemaleAge: 79,MinFemaleAge: 20,MaxMaleAge: 80,MinMaleAge: 21}}).Error)
	// topup
	agencyTopup := model.AgencyTopUp{BlueAgencytopup: gen.BlueAgencytopup{AgencyID: agency.ID}}
	res := ctx.DB().FirstOrCreate(&agencyTopup,"agency_id = ?",agency.ID)
	assert.NoError(t, res.Error)
	if res.RowsAffected == 0 {

	}*/
	// role
	role := model.Role{BlueRole: gen.BlueRole{Name: "管理员", Category: 2, AgencyID: base.NullIDFromID(agency.ID)}}
	assert.NoError(t, ctx.DB().Create(role).Error)
	user := model.User{BlueUser: gen.BlueUser{Mobile: loginMobile}}
	userAgency := model.UserAgency{BlueUseragency: gen.BlueUseragency{Name: "沈佩", UserID: base.NullIDFromID(user.ID), AgencyID: agency.ID, RoleID: base.NullIDFromID(role.ID)}}
	// todo create token(Auth) 这里需要修改一下创建token的方法 参考:func (h BlueUserHandler) Login(ctx *app.Context)
}
