package context

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Context holds the request and response writer and provides utility methods.
type Context struct {
	Req *http.Request
	Res http.ResponseWriter
}

// NewContext creates a new Context instance.
//
// This function initializes a new Context with the provided request and response writer.
//
// Example usage:
//
//	ctx := context.NewContext(req, res)
func NewContext(req *http.Request, res http.ResponseWriter) *Context {
	return &Context{Req: req, Res: res}
}

// GetJSONBody retrieves the parsed JSON body from the request context.
func (c *Context) GetJSONBody() (map[string]interface{}, bool) {
	key := "jsonBody"
	jsonBody, ok := c.Req.Context().Value(key).(map[string]interface{})
	return jsonBody, ok
}

// JSON sends a JSON response with the given status code.
//
// This method sets the Content-Type to application/json and writes the provided value as a JSON response.
//
// Parameters:
//
//	status (int): The HTTP status code to send with the response.
//	v (interface{}): The data to encode as JSON and send in the response.
//
// Example usage:
//
//	ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
func (c *Context) JSON(status int, v interface{}) {
	c.Res.Header().Set("Content-Type", "application/json")
	c.Res.WriteHeader(status)
	json.NewEncoder(c.Res).Encode(v)
}

// Error sends an error response with the given status code and message.
//
// This method sets the Content-Type to application/json and writes an error message with the specified HTTP status code.
//
// Parameters:
//
//	status (int): The HTTP status code to send with the response.
//	message (string): The error message to include in the response.
//
// Example usage:
//
//	ctx.Error(http.StatusBadRequest, "Invalid request")
func (c *Context) Error(status int, message string) {
	c.Res.Header().Set("Content-Type", "application/json")
	c.Res.WriteHeader(status)
	json.NewEncoder(c.Res).Encode(map[string]string{"error": message})
}

// Body parses the JSON request body into the provided interface.
//
// This method decodes the JSON body of the request into the provided value.
//
// Parameters:
//
//	v (interface{}): The value to decode the JSON into.
//
// Returns:
//
//	error: An error if JSON decoding fails.
//
// Example usage:
//
//	var data map[string]interface{}
//	err := ctx.Body(&data)
func (c *Context) Body(v interface{}) error {
	if c.Req.Body == nil {
		return errors.New("request body is nil")
	}
	bodyBytes, err := io.ReadAll(c.Req.Body)
	if err != nil {
		return err
	}
	if len(bodyBytes) == 0 {
		return errors.New("empty request body")
	}
	// Reset the body so it can be read again later if needed
	c.Req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(v)
}

// Redirect sends a redirect response to the given URL.
//
// This method sends an HTTP redirect to the specified URL with the provided status code.
//
// Parameters:
//
//	status (int): The HTTP status code for the redirect (e.g., http.StatusFound).
//	url (string): The URL to redirect to.
//
// Example usage:
//
//	ctx.Redirect(http.StatusSeeOther, "/new-url")
func (c *Context) Redirect(status int, url string) {
	http.Redirect(c.Res, c.Req, url, status)
}

// SetCookie adds a cookie to the response.
//
// This method sets a cookie with the given name, value, and options.
//
// Parameters:
//
//	name (string): The name of the cookie.
//	value (string): The value of the cookie.
//	maxAge (int): The maximum age of the cookie in seconds.
//	path (string): The path for which the cookie is valid.
//
// Example usage:
//
//	ctx.SetCookie("session_id", "123456", 3600, "/")
func (c *Context) SetCookie(name, value string, maxAge int, path string) {
	http.SetCookie(c.Res, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		HttpOnly: true,
	})
}

// GetCookie retrieves a cookie value from the request.
//
// This method fetches the value of a cookie with the specified name from the request.
//
// Parameters:
//
//	name (string): The name of the cookie to retrieve.
//
// Returns:
//
//	(string, bool): The value of the cookie and a boolean indicating if the cookie was found.
//
// Example usage:
//
//	value, ok := ctx.GetCookie("session_id")
func (c *Context) GetCookie(name string) (string, bool) {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}
