package test

import (
	"io/ioutil"
	"lessgo/internal/core/controller"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type TestController struct {
	controller.BaseController
	TestService TestServiceInterface
}

func NewTestController(service TestServiceInterface) *TestController {
	return &TestController{
		TestService: service,
	}
}

func (tc *TestController) RegisterRoutes(mux *mux.Router) {
	mux.HandleFunc("/example", tc.handleTest).Methods("GET")
	mux.HandleFunc("/about", tc.handleAbout).Methods("GET")
}

func (tc *TestController) handleTest(w http.ResponseWriter, r *http.Request) {
	response := tc.TestService.DoSomething()
	w.Write([]byte(response))
}

func (tc *TestController) handleAbout(w http.ResponseWriter, r *http.Request) {
	// Get the absolute path of the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Unable to get current directory", http.StatusInternalServerError)
		return
	}
	log.Print(cwd)
	// Determine the path to the README file
	readmePath := filepath.Join(cwd, "README.md")

	// Read the README file
	data, err := ioutil.ReadFile(readmePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set content type and write the file content
	w.Header().Set("Content-Type", "text/markdown")
	w.Write(data)
}
