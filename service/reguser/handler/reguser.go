package handler

import (
	"context"
	reguser "ihome/service/reguser/proto"
	"ihome/web/dao/mysql"
	"ihome/web/dao/redis"
)

type Reguser struct{}

func (r Reguser) RegisterUser(ctx context.Context, request *reguser.RegisterRequest, response *reguser.RegisterResponse) error {
	// 校验sms
	if err := redis.CheckSms(request.Phone, request.SmsInput); err != nil {
		response.Errmsg = "redis checksms fail: " + err.Error()
		return err
	}

	// 数据存入数据库
	if err := mysql.RegisterNewAccount(request.Phone, request.Password); err != nil {
		response.Errmsg = "mysql register new account fail: " + err.Error()
		return err
	}
	return nil
}

// Return a new handler
func New() *Reguser {
	return &Reguser{}
}
