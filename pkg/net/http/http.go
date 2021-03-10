package http

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
	"net/http"
	"reflect"
)

func GRPCProxyWrapper(h interface{}) echo.HandlerFunc {
	v := reflect.ValueOf(h)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic("reflect error: handler must be func")
	}
	return func(c echo.Context) error {
		var req = reflect.New(t.In(1).Elem()).Interface()
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"msg":  err.Error(),
				"code": 40023,
			})
		}
		var md = metadata.MD{}
		for k, vs := range c.Request().Header {
			for _, v := range vs {
				bs := bytes.TrimFunc([]byte(v), func(r rune) bool {
					return r == '\n' || r == '\r' || r == '\000'
				})
				md.Append(k, string(bs))
			}
		}
		ctx := metadata.NewOutgoingContext(c.Request().Context(), md)
		values := v.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)})
		resp, err := values[0], values[1]
		if !err.IsNil() {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"msg":  err.Interface(),
				"code": 50000,
			})
		}
		return c.JSON(http.StatusOK, resp.Interface())
	}
}

