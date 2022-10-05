package mysql

import (
	"errors"
	"fmt"
	"ihome/web/model"
	"ihome/web/utils"
)

func RegisterNewAccount(phone, password string) error {
	db, err := Get(IHOME)
	if err != nil {
		return err
	}
	// 1. 检测用户是否已注册
	rows, err := db.Table("users").Where("mobile = ?", phone).Rows()
	if err != nil {
		return err
	}
	if rows.Next() {
		return errors.New("phone was already registered")
	}
	// 2. 持久化一个新用户信息
	newOne := model.User{}
	newOne.Mobile = phone
	pwd_hash := utils.Pwd2Hash(password)
	newOne.Password_hash = pwd_hash
	db.Create(&newOne)
	return nil
}

func GetArea() ([]model.Area, error) {
	var areas []model.Area
	db, err := Get(IHOME)
	if err != nil {
		return nil, err
	}
	db.Find(&areas)
	return areas, nil
}

func CheckUser(phone, password string) (name string, err error) {
	db, err := Get(IHOME)
	if err != nil {
		return "", errors.New("Get db," + err.Error())
	}
	pwd_hash := utils.Pwd2Hash(password)
	var user model.User
	db.Where("mobile = ?", phone).Find(&user)
	if user.Mobile == "" {
		return "", errors.New("no such user")
	}
	if user.Password_hash != pwd_hash {
		fmt.Println("stored password: ", user.Password_hash)
		fmt.Println("input password : ", pwd_hash)
		return "", errors.New("password incorrect")
	}
	name = user.Name
	return
}

func GetUserInfoByName(userName string) (model.User, error) {
	var user model.User
	db, err := Get(IHOME)
	if err != nil {
		return model.User{}, errors.New("Get db," + err.Error())
	}

	db.Omit("password_hash").Where("name = ?", userName).First(&user)
	return user, nil
}

func GetUserInfoById(userId string) (model.User, error) {
	var user model.User
	db, err := Get(IHOME)
	if err != nil {
		return model.User{}, errors.New("Get db," + err.Error())
	}

	db.Omit("password_hash").Where("id = ?", userId).First(&user)
	return user, nil
}

func PutUserName(oldName, newName string) error {
	db, err := Get(IHOME)
	if err != nil {
		return errors.New("Get db," + err.Error())
	}

	// 更新名字
	return db.Model(model.User{}).Where("name = ?", oldName).Update("name", newName).Error
}

func SaveAvatarURL(name, avatarURL string) error {
	db, err := Get(IHOME)
	if err != nil {
		return errors.New("Get db," + err.Error())
	}
	return db.Model(model.User{}).Where("name = ?", name).Update("avatar_url", avatarURL).Error
}

func SaveRealName(name, realName, id string) error {
	db, err := Get(IHOME)
	if err != nil {
		return errors.New("Get db," + err.Error())
	}
	return db.Model(model.User{}).Where("name = ?", name).Updates(model.User{Real_name: realName, Id_card: id}).Error
}

func GetUserAuth(name string) (model.User, error) {
	db, err := Get(IHOME)
	if err != nil {
		return model.User{}, errors.New("Get db," + err.Error())
	}
	var u model.User
	err = db.Where("name = ?", name).First(&u).Error
	return u, err
}

func GetUserAvatarURL(userID string) (string, error) {
	u, err := GetUserInfoById(userID)
	if err != nil {
		return "", err
	}
	return u.Avatar_url, nil
}
