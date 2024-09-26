package src

type IRootService interface{}

type RootService struct {
	// Add any shared dependencies or methods here
}

func NewRootService() *RootService {
	return &RootService{}
}

// Add methods that interact with other services
