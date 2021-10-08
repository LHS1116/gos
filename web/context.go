package web

import (
	"encoding/json"
	"net/http"
)

type GosContext struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Status  int
}
type H map[string]interface{}

func newContext(r *http.Request, w http.ResponseWriter) *GosContext {
	return &GosContext{r, w, 0}
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

func (c *GosContext) SetHeader(k, v string) {
	c.Writer.Header().Set(k, v)
}

func (c *GosContext) SetStatus(code int) {
	c.Writer.WriteHeader(code)
	c.Status = code
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
