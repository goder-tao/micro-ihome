package user

type RegistryPayload struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	SmsCode  string `json:"sms_code"`
}

type LoginPayLoad struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}
