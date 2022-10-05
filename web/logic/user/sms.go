package user

import (
	"context"
	"errors"
	"github.com/micro/micro/v3/service"
	Nsms "ihome/service/sms/proto"
	"log"
)

// SendSms 向指定号码发送短信验证码并且返回发送的验证码内容
// 仅利用随机数模拟逻辑，受限于短信服务的签名认证
// 由微服务实现短信验证码的发送和存储到redis
func SendSms(phone string) error {
	svc := service.New()
	cli := Nsms.NewSmsService("sms", svc.Client())
	rsp, err := cli.SendSms(context.Background(), &Nsms.SmsRequest{Phone: phone})
	if err != nil {
		log.Println("SendSms err:", err)
		return err
	}
	if rsp.Errmsg != "" {
		log.Println("remote call errmsg:", rsp.Errmsg)
		return errors.New(rsp.Errmsg)
	}
	return nil
}
