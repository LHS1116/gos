package web

import (
	"encoding/json"
	"net/http"
)

type GosContext struct {
	Request     *http.Request
	Writer      http.ResponseWriter
	Status      int
	index       int
	middlewares []Middleware
	pathParams  map[string]string
}
type H map[string]interface{}

func newContext(r *http.Request, w http.ResponseWriter, middlewares []Middleware) *GosContext {
	return &GosContext{
		Request:     r,
		Writer:      w,
		index:       -1,
		middlewares: middlewares,
	}
}

func (c *GosContext) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *GosContext) DefaultQuery(key string, defaultValue string) string {
	values := c.Request.URL.Query()
	if values == nil {
		return defaultValue
	}
	v, ok := values[key]
	if !ok {
		return defaultValue
	}
	return v[0]
}

func (c *GosContext) PathParam(key string, defaultValue string) string {
	if c.pathParams != nil {
		v, ok := c.pathParams[key]
		if ok {
			return v
		}
	}
	return defaultValue
}

func (c *GosContext) FullPath() string {
	return c.Request.Host + c.Request.RequestURI
}

func (c *GosContext) SetHeader(k, v string) {
	c.Writer.Header().Set(k, v)
}

func (c *GosContext) SetStatus(code int) {
	c.Writer.WriteHeader(code)
	c.Status = code
}

func (c *GosContext) Next() {
	c.index++
	for c.index < len(c.middlewares) {
		c.middlewares[c.index]()(c)
		c.index++
	}
}

func (c *GosContext) JSON(code int, data interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *GosContext) String(code int, str string) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	_, err := c.Writer.Write([]byte(str))
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *GosContext) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	_, err := c.Writer.Write([]byte(html))
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
