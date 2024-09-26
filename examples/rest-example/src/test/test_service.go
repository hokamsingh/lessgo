package test

import (
	"log"

	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type ITestService interface{}

type TestService struct {
	LessGo.BaseService
}

func NewTestService() *TestService {
	return &TestService{}
}

func (es *TestService) DoSomething() string {
	log.Print("Service Logic Executed")
	return "Service Logic Executed"
}
