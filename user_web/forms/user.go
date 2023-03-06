package forms

//表单提交格式
type PassWordLoginForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required"` //一些格式验证得自己写函数
	PassWord string `form:"password" json:"passWord" binding:"required,min=3,max=18"`
}
