/*
Package service provides a base structure and interface for defining and implementing services in the application.

This package defines a Service interface, which can be extended to implement various services. Additionally, it provides a BaseService struct with a default method implementation.
*/
package service

// Service defines the interface for all services in the application.
// Implementations of this interface can provide specific functionalities
// required by different parts of the application.
type Service interface {
}

// BaseService provides a default implementation of the Service interface.
// This struct can be embedded in other service implementations to inherit
// common functionalities or to be extended with custom methods.
type BaseService struct{}

// PerformTask is a method of BaseService that performs a generic task.
// This method can be overridden by services embedding BaseService to provide
// specific behavior or functionality.
//
// Example:
//
//	type MyService struct {
//		service.BaseService
//	}
//
//	func (s *MyService) PerformTask() {
//		// Custom task implementation
//		fmt.Println("Performing a custom task")
//	}
//
//	func main() {
//		s := MyService{}
//		s.PerformTask() // Outputs: Performing a custom task
//	}
func (bs *BaseService) PerformTask() {

}
