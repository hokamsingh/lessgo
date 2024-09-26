package user

import (
	"log"

	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type IUserService interface{}

type UserService struct {
	LessGo.BaseService
}

func NewUserService() *UserService {
	return &UserService{}
}

func (es *UserService) DoSomething() string {
	log.Print("Service Logic Executed")
	return "Service Logic Executed"
}
