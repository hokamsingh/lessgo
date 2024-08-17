package test

import (
	"log"

	"github.com/hokamsingh/lessgo/internal/core/service"
)

type TestServiceInterface interface {
	DoSomething() string
}

type TestService struct {
	service.BaseService
}

func NewTestService() *TestService {
	return &TestService{}
}

func (es *TestService) DoSomething() string {
	log.Print("Service Logic Executed")
	return "Service Logic Executed"
}
