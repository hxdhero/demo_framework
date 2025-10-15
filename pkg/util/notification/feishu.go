package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"lls_api/pkg/config"
	"lls_api/pkg/rerr"
	"time"
)

var larkClient *lark.Client

func InitLarkClient() error {
	larkClient = lark.NewClient(config.C.Notification.Feishu.AppID, config.C.Notification.Feishu.Secret)
	_, err := larkClient.Verification.V1.Verification.Get(context.Background())
	if err != nil {
		return rerr.Wrap(err)
	}
	return nil
}

// SendFeishuUserText 给单个用户发送文本
func SendFeishuUserText(userId, msg string) error {
	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`user_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(userId).
			MsgType(`text`).
			Content(fmt.Sprintf(`{"text":"%s"}`, msg)).
			Build()).
		Build()

	// 发起请求
	resp, err := larkClient.Im.V1.Message.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		return rerr.Wrap(err)
	}

	// 服务端错误处理
	if !resp.Success() {
		return rerr.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}

	return nil
}

// SendFeishuChatCard 给群组发送文本
func SendFeishuChatCard(chartId, msg string) error {
	templateIDStr := config.C.Notification.Feishu.TemplateID
	type TemplateVariable struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}
	type Data struct {
		TemplateID       string           `json:"template_id"`
		TemplateVariable TemplateVariable `json:"template_variable"`
	}
	type Content struct {
		Type string `json:"type"`
		Data Data   `json:"data"`
	}
	bs, err := json.Marshal(Content{
		Type: "template",
		Data: Data{
			TemplateID: templateIDStr,
			TemplateVariable: TemplateVariable{
				Title: "立利顺",
				Text:  fmt.Sprintf("%s\n%s", time.Now().Format(time.DateTime), msg),
			},
		},
	})
	if err != nil {
		return rerr.Wrap(err)
	}
	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`chat_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(chartId).
			MsgType(`interactive`).
			Content(string(bs)).
			Build()).
		Build()
	// 发起请求
	resp, err := larkClient.Im.V1.Message.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		return rerr.Wrap(err)
	}

	// 服务端错误处理
	if !resp.Success() {
		return rerr.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}

	return nil
}
