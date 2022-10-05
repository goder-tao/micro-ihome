package handler

import (
	"context"
	"fmt"
	"ihome/web/dao/redis"
	"log"
	"math/rand"
	"time"

	sms "ihome/service/sms/proto"
)

type Sms struct{}

// Return a new handler
func New() *Sms {
	return &Sms{}
}

func (e *Sms) SendSms(ctx context.Context, req *sms.SmsRequest, rsp *sms.SmsResponse) error {
	// 模拟短信验证码发送
	randn := rand.New(rand.NewSource(time.Now().UnixNano()))
	smsCode := fmt.Sprintf("%06d", randn.Int()%1000000)
	log.Println("sms code: ", smsCode)

	// 存入smscode存入redis
	if err := redis.SaveSms(req.Phone, smsCode); err != nil {
		rsp.Errmsg = err.Error()
	}
	
	return nil
}
