package forms

//表单提交格式
type PassWordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile_validate"` //一些格式验证得自己写函数
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=18"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile_validate"` //手机号码格式有规范可寻， 自定义validator
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	//Code string `form:"code" json:"code" binding:"required,min=6,max=6"` //验证码暂时不需要
}
