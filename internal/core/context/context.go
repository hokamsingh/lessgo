package context

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Context holds the request and response writer and provides utility methods.
type Context struct {
	Req          *http.Request
	Res          http.ResponseWriter
	responseSent bool // Track whether a response has been sent
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
	if c.responseSent {
		log.Fatal("Response already sent")
		return
	}
	c.Res.Header().Set("Content-Type", "application/json")
	c.Res.WriteHeader(status)
	json.NewEncoder(c.Res).Encode(v)
	c.responseSent = true
	c.Res.(http.Flusher).Flush() // Ensures the data is sent to the client
}

// Send sends a plain text response.
//
// This method sets the Content-Type to text/plain and writes the provided value as a response.
//
// Parameters:
//
//	v (any): The value to send in the response. It will be converted to a string.
//
// Example usage:
//
//	ctx.Send("Hello, World!")
func (c *Context) Send(v any) {
	if c.responseSent {
		log.Fatal("Response already sent")
		return
	}
	c.SetHeader("Content-Type", "text/plain")
	c.Res.WriteHeader(http.StatusOK)
	c.Res.Write([]byte(fmt.Sprint(v)))
	c.responseSent = true
	c.Res.(http.Flusher).Flush() // Ensures the data is sent to the client
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
	if c.responseSent {
		log.Fatal("Response already sent")
		return
	}
	c.Res.Header().Set("Content-Type", "application/json")
	c.Res.WriteHeader(status)
	json.NewEncoder(c.Res).Encode(map[string]string{"error": message})
	// Close the response after sending the error
	c.responseSent = true
	c.Res.(http.Flusher).Flush() // Ensures the data is sent to the client
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
	if c.responseSent {
		log.Fatal("Response already sent")
		return
	}
	http.Redirect(c.Res, c.Req, url, status)
}

type SameSite int

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

// SetCookie adds a cookie to the response.
//
// This method sets a cookie with the given name, value, and options.
//
// Parameters:
//
//			name (string): The name of the cookie.
//			value (string): The value of the cookie.
//			maxAge (int): The maximum age of the cookie in seconds.
//			path (string): The path for which the cookie is valid.
//		    httpOnly(bool): make it http only cookie.(can only be sent from http request)
//		    secure(bool): ensures that the cookie is only sent over HTTPS connections, providing protection against man-in-the-middle (MITM) attacks
//			sameSite(http.SameSite): Strict(SameSiteStrictMode): Most restrictive, no cross-site requests.
//	     Lax(SameSiteLaxMode): Allows cookies to be sent with top-level navigations, but not with other cross-site requests.
//
//			None(SameSiteNoneMode): No restrictions on sending cookies with cross-site requests, but must be used with Secure.
//
//			Default(SameSiteDefaultMode): Equivalent to Lax
//
// Example usage:
//
//	ctx.SetCookie("auth_token", "0xc000013a", 60, "", true, false, http.SameSiteDefaultMode)
func (c *Context) SetCookie(name, value string, maxAge int, path string, httpOnly bool, secure bool, sameSite http.SameSite) {
	http.SetCookie(c.Res, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: sameSite,
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

// GetParam retrieves a URL parameter from the request.
// Assumes that parameters are stored in the request context.
func (c *Context) GetParam(name string) (string, bool) {
	params, ok := c.GetAllParams()
	if !ok {
		return "", false
	}
	value, found := params[name]
	return value, found
}

// GetAllParams retrieves all URL parameters from the request.
//
// This method returns a map containing all URL parameters with their respective values
// and a boolean indicating whether any parameters were found.
//
// Returns:
//
//	(map[string]string, bool): A map of URL parameters and a boolean indicating if any were found.
func (c *Context) GetAllParams() (map[string]string, bool) {
	params := mux.Vars(c.Req)
	if len(params) == 0 {
		return nil, false
	}
	return params, true
}

// GetQuery retrieves a query parameter from the request URL.
func (c *Context) GetQuery(name string) (string, bool) {
	values := c.Req.URL.Query()
	value := values.Get(name)
	return value, value != ""
}

// GetAllQuery retrieves all query parameters as a JSON object.
func (c *Context) GetAllQuery() (map[string]interface{}, error) {
	queryMap := make(map[string]interface{})
	for key, values := range c.Req.URL.Query() {
		if len(values) > 1 {
			queryMap[key] = values
		} else {
			queryMap[key] = values[0]
		}
	}
	return queryMap, nil
}

// GetHeader retrieves a header value from the request.
func (c *Context) GetHeader(name string) string {
	return c.Req.Header.Get(name)
}

// SetHeader sets a header value for the response.
func (c *Context) SetHeader(name, value string) {
	c.Res.Header().Set(name, value)
}
