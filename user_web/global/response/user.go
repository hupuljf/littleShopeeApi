package response

//用于给web的http调用返回json数据
type UserResponse struct {
	Id       int32  `json:"id"`
	NickName string `json:"name"`
	Gender   string `json:"gender"`
	Mobile   string `json:"mobile"`
	Birthday string `json:"birthday"`
}
