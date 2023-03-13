package models

import "github.com/dgrijalva/jwt-go"

//使用jwt 自定义的palyload id 昵称 角色
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}
