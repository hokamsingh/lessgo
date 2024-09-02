
---

## Package `context`

The `context` package provides utilities to handle HTTP requests and responses effectively.

### Types

#### `Context`

A `Context` holds the HTTP request and response writer, and provides methods to manage request and response operations.

### Methods

#### `NewContext`

```go
func NewContext(req *http.Request, res http.ResponseWriter) *Context
```

Creates a new `Context` instance with the given request and response writer.

**Usage:**

```go
ctx := context.NewContext(req, res)
```

#### `GetJSONBody`

```go
func (c *Context) GetJSONBody() (map[string]interface{}, bool)
```

Retrieves the parsed JSON body from the request context.

**Returns:**

- The JSON body as a map and a boolean indicating if the body was found.

**Usage:**

```go
jsonBody, ok := ctx.GetJSONBody()
```

#### `JSON`

```go
func (c *Context) JSON(status int, v interface{})
```

Sends a JSON response with the given status code.

**Parameters:**

- `status`: HTTP status code.
- `v`: Data to encode as JSON.

**Usage:**

```go
ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
```

#### `Send`

```go
func (c *Context) Send(v any)
```

Sends a plain text response.

**Parameters:**

- `v`: Value to send in the response.

**Usage:**

```go
ctx.Send("Hello, World!")
```

#### `Error`

```go
func (c *Context) Error(status int, message string)
```

Sends an error response with the given status code and message.

**Parameters:**

- `status`: HTTP status code.
- `message`: Error message.

**Usage:**

```go
ctx.Error(http.StatusBadRequest, "Invalid request")
```

#### `Body`

```go
func (c *Context) Body(v interface{}) error
```

Parses the JSON request body into the provided interface.

**Parameters:**

- `v`: Value to decode the JSON into.

**Returns:**

- An error if JSON decoding fails.

**Usage:**

```go
var data map[string]interface{}
err := ctx.Body(&data)
```

#### `Redirect`

```go
func (c *Context) Redirect(status int, url string)
```

Sends a redirect response to the given URL.

**Parameters:**

- `status`: HTTP status code for the redirect.
- `url`: URL to redirect to.

**Usage:**

```go
ctx.Redirect(http.StatusSeeOther, "/new-url")
```

#### `SetCookie`

```go
func (c *Context) SetCookie(name, value string, maxAge int, path string, httpOnly bool, secure bool, sameSite http.SameSite)
```

Adds a cookie to the response.

**Parameters:**

- `name`: Name of the cookie.
- `value`: Value of the cookie.
- `maxAge`: Maximum age of the cookie in seconds.
- `path`: Path for which the cookie is valid.
- `httpOnly`: If true, the cookie is accessible only via HTTP(S).
- `secure`: If true, the cookie is sent only over HTTPS.
- `sameSite`: SameSite attribute for the cookie.

**Usage:**

```go
ctx.SetCookie("auth_token", "0xc000013a", 60, "/", true, true, http.SameSiteLax)
```

#### `GetCookie`

```go
func (c *Context) GetCookie(name string) (string, bool)
```

Retrieves a cookie value from the request.

**Parameters:**

- `name`: Name of the cookie to retrieve.

**Returns:**

- Value of the cookie and a boolean indicating if the cookie was found.

**Usage:**

```go
value, ok := ctx.GetCookie("session_id")
```

#### `GetParam`

```go
func (c *Context) GetParam(name string) (string, bool)
```

Retrieves a URL parameter from the request.

**Parameters:**

- `name`: Name of the parameter to retrieve.

**Returns:**

- Value of the parameter and a boolean indicating if the parameter was found.

**Usage:**

```go
value, ok := ctx.GetParam("id")
```

#### `GetAllParams`

```go
func (c *Context) GetAllParams() (map[string]string, bool)
```

Retrieves all URL parameters from the request.

**Returns:**

- A map of URL parameters and a boolean indicating if any parameters were found.

**Usage:**

```go
params, ok := ctx.GetAllParams()
```

#### `GetQuery`

```go
func (c *Context) GetQuery(name string) (string, bool)
```

Retrieves a query parameter from the request URL.

**Parameters:**

- `name`: Name of the query parameter to retrieve.

**Returns:**

- Value of the query parameter and a boolean indicating if the parameter was found.

**Usage:**

```go
value, ok := ctx.GetQuery("search")
```

#### `GetAllQuery`

```go
func (c *Context) GetAllQuery() (map[string]interface{}, error)
```

Retrieves all query parameters as a map.

**Returns:**

- A map of query parameters and their values, and an error if the parameters could not be retrieved.

**Usage:**

```go
queryParams, err := ctx.GetAllQuery()
```

#### `GetHeader`

```go
func (c *Context) GetHeader(name string) string
```

Retrieves a header value from the request.

**Parameters:**

- `name`: Name of the header to retrieve.

**Returns:**

- Value of the header.

**Usage:**

```go
value := ctx.GetHeader("Authorization")
```

#### `SetHeader`

```go
func (c *Context) SetHeader(name, value string)
```

Sets a header value for the response.

**Parameters:**

- `name`: Name of the header to set.
- `value`: Value of the header.

**Usage:**

```go
ctx.SetHeader("X-Custom-Header", "value")
```

#### `Status`

```go
func (c *Context) Status(code int)
```

Sets the HTTP response code.

**Parameters:**

- `code`: HTTP status code.

**Usage:**

```go
ctx.Status(http.StatusNotFound)
```

#### `FileAttachment`

```go
func (c *Context) FileAttachment(filepath, filename string)
```

Sends a file as an attachment for download.

**Parameters:**

- `filepath`: Path to the file.
- `filename`: Name of the file as it will appear to the client.

**Usage:**

```go
ctx.FileAttachment("/path/to/file.txt", "file.txt")
```

---