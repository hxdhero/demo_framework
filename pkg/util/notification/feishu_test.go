package notification

import (
	"fmt"
	"lls_api/pkg/config"
	"testing"
)

func TestSendFeishuUserText(t *testing.T) {
	test := false
	if !test {
		return
	}
	config.InitConfig()
	if err := InitLarkClient(); err != nil {
		panic(err)
	}
	if err := SendFeishuUserText("f94b5bb5", "hello"); err != nil {
		fmt.Println(err)
	}
}

func TestSendFeishuChatText(t *testing.T) {
	test := false
	if !test {
		return
	}
	config.InitConfig()
	if err := InitLarkClient(); err != nil {
		panic(err)
	}
	msg := "这是测试信息,有一个地方出错了"
	if err := SendFeishuChatCard(config.C.Notification.Feishu.ErrChatID, msg); err != nil {
		fmt.Println(err)
	}
}
