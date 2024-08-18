package test

import (
	"log"

	core "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type TestServiceInterface interface {
	DoSomething() string
}

type TestService struct {
	core.Service
}

func NewTestService() *TestService {
	return &TestService{}
}

func (es *TestService) DoSomething() string {
	log.Print("Service Logic Executed")
	return "Service Logic Executed"
}
