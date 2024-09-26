package test

import (
	"net/http"

	LessGo "github.com/hokamsingh/lessgo/pkg/lessgo"
)

type TestController struct {
	Path    string
	Service TestService
}

func NewTestController(service *TestService, path string) *TestController {
	// if !LessGo.ValidatePath(path) {
	// 	log.Fatalf("Invalid path provided: %s", path)
	// }
	return &TestController{
		Service: *service,
		Path:    path,
	}
}

func (tc *TestController) RegisterRoutes(r *LessGo.Router) {
	tr := r.SubRouter(tc.Path)
	tr.Get("/ping", func(ctx *LessGo.Context) {
		// ctx.JSON(200, map[string]string{"message": "pong"})
		ctx.Send("pong")
	})

	tr.Get("/info", func(ctx *LessGo.Context) {
		info := tc.Service.DoSomething()
		ctx.Send(info)
	})

	tr.Get("/user/{id}", func(ctx *LessGo.Context) {
		// Get all URL params
		params, ok := ctx.GetAllParams()
		id := params["id"]
		if !ok {
			ctx.Error(400, "no params found")
			return
		}
		// Get all query params as JSON
		queryParams, _ := ctx.GetAllQuery()
		// Set a custom header
		ctx.SetHeader("X-Custom-Header", "MyValue")
		cookie, ok := ctx.GetCookie("auth_token")
		if !ok {
			// ctx.Error(400, "Bad Request")
			ctx.SetCookie("auth_token", "0xc000013a", 60, "", true, false, http.SameSiteDefaultMode)
		}
		ctx.JSON(200, map[string]interface{}{
			"params":      params,
			"queryParams": queryParams,
			"id":          id,
			"cookie":      cookie,
		})
	})

	tr.Post("/submit", func(ctx *LessGo.Context) {
		var body User
		ctx.Body(&body)
		ctx.JSON(200, body)
	})

	tr.Delete("/{id}", func(ctx *LessGo.Context) {
		var id string
		id, _ = ctx.GetParam("id")
		ctx.Error(400, id)
	})
}

// func (tc *TestController) handleTest(w http.ResponseWriter, r *http.Request) {
// 	response := tc.TestService.DoSomething()
// 	w.Write([]byte(response))
// }

// func (tc *TestController) handleAbout(w http.ResponseWriter, r *http.Request) {
// 	// Get the absolute path of the current working directory
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		http.Error(w, "Unable to get current directory", http.StatusInternalServerError)
// 		return
// 	}
// 	log.Print(cwd)
// 	// Determine the path to the README file
// 	readmePath := filepath.Join(cwd, "README.md")

// 	// Read the README file
// 	data, err := ioutil.ReadFile(readmePath)
// 	if err != nil {
// 		http.Error(w, "File not found", http.StatusNotFound)
// 		return
// 	}

// 	// Set content type and write the file content
// 	w.Header().Set("Content-Type", "text/markdown")
// 	w.Write(data)
// }

// func (c *TestController) handleData(w http.ResponseWriter, r *http.Request) {
// 	// Retrieve the parsed JSON from the request context
// 	parsedData, ok := r.Context().Value("jsonBody").(map[string]interface{})
// 	if !ok {
// 		http.Error(w, "Failed to retrieve JSON data", http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the response header to JSON
// 	w.Header().Set("Content-Type", "application/json")

// 	// Encode the parsed data back into JSON and write to the response
// 	err := json.NewEncoder(w).Encode(parsedData)
// 	if err != nil {
// 		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
// 	}
// }
