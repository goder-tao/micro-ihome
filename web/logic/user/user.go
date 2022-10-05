package user

import (
	"context"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/micro/micro/v3/service"
	avatar "ihome/service/avatar/proto"
	realname "ihome/service/realname/proto"
	"ihome/web/dao/mysql"
	"ihome/web/dao/redis"
	"ihome/web/model"
	"ihome/web/utils"
)

const (
	SESSION_USER_NAME = "USER_NAME"
)

func Register(phone, inputSms, password string) error {
	//svc := service.New()
	//cli := reguser.NewReguserService("reguser", svc.Client())
	//_, err := cli.RegisterUser(context.Background(), &reguser.RegisterRequest{Phone: phone, Password: password, SmsInput: inputSms})
	//if err != nil {
	//	return err
	//}
	//return nil
	// 校验sms
	if err := redis.CheckSms(phone, inputSms); err != nil {
		return errors.New("redis checksms fail: " + err.Error())
	}

	// 数据存入数据库
	if err := mysql.RegisterNewAccount(phone, password); err != nil {
		return errors.New("mysql register new account fail: " + err.Error())
	}
	return nil
}

func CheckUser(phone, password string) (name string, err error) {
	return mysql.CheckUser(phone, password)
}

func SaveSessionName(c *gin.Context, userName string) error {
	s := sessions.Default(c)
	s.Set(SESSION_USER_NAME, userName)
	err := s.Save()
	return err
}

func GetSessionName(c *gin.Context) (string, error) {
	s := sessions.Default(c)
	name := s.Get(SESSION_USER_NAME)
	if name == nil {
		return "", errors.New(utils.RecodeText(utils.RECODE_SESSIONERR))
	}
	return name.(string), nil
}

func DeleteSession(c *gin.Context) error {
	s := sessions.Default(c)
	s.Delete(SESSION_USER_NAME)
	return s.Save()
}

func GetUserInfo(userName string) (model.User, error) {
	return mysql.GetUserInfoByName(userName)
}

func PutUserName(c *gin.Context, newName string) error {
	// 获取旧名
	oldName, err := GetSessionName(c)
	if err != nil {
		return err
	}
	// 更新数据库中用户名
	if err := mysql.PutUserName(oldName, newName); err != nil {
		return err
	}
	// 更新session
	return SaveSessionName(c, newName)
}

// SaveAvatarByte 存储头像的[]byte到fdfs
func SaveAvatarByte(buf []byte, extName string) (string, error) {
	svc := service.New()
	cli := avatar.NewAvatarService("avatar", svc.Client())
	resp, err := cli.SaveBuffer(context.Background(), &avatar.AvatarRequest{Buf: buf, ExeName: extName})
	if err != nil {
		return "", err
	}
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}
	return resp.FileID, nil
}

// SaveAvatarURL 保存avatarURL
func SaveAvatarURL(c *gin.Context, avatarURL string) error {
	// 1. 从session获取用户名
	name, err := GetSessionName(c)
	if err != nil {
		return err
	}

	// 2. 存入mysql
	return mysql.SaveAvatarURL(name, avatarURL)
	// 3.可能会失败，所以应该保留旧的avatarURL，以防保存失败，保存成功的时候应该删除掉旧的数据
}

// AuthUser 实名认证
func AuthUser(realName, id string, c *gin.Context) error {
	// 认证服务认证
	svc := service.New()
	cli := realname.NewRealnameService("realname", svc.Client())
	rsp, err := cli.AuthRealName(context.Background(), &realname.RealNameRequest{
		Realname: realName,
		IdNumber: id,
	})
	if err != nil {
		return err
	}
	if rsp.Errmsg != "" {
		return errors.New(rsp.Errmsg)
	}

	// 持久化认证结果
	name, err := GetSessionName(c)
	if err != nil {
		return err
	}
	return mysql.SaveRealName(name, realName, id)
}

func GetUserAuth(c *gin.Context) (model.User, error) {
	name, err := GetSessionName(c)
	if err != nil {
		return model.User{}, err
	}
	return mysql.GetUserAuth(name)
}
