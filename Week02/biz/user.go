package biz

import (
	"github.com/crosserclaws/Go-000/Week02/dao"
)

// FetchUserInfo is an example biz (i.e. service) layer function.
func FetchUserInfo(id int) (*dao.User, error) {
	return dao.QueryUserByID(id)
}

// ListUsers is an example biz (i.e. service) layer function.
func ListUsers() ([]*dao.User, error) {
	return dao.QueryAllUsers()
}
