package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/assert"
)

func TestRouteRawPath(t *testing.T) {
	route := gin.New()
	route.UseRawPath = true

	route.POST("/project/:name/build/:num", func(c *gin.Context) {
		name := c.Params.ByName("name")
		num := c.Params.ByName("num")

		assert.Equal(t, name, c.Param("name"))
		assert.Equal(t, num, c.Param("num"))

		assert.Equal(t, "Some/Other/Project", name)
		assert.Equal(t, "222", num)
	})

	w := PerformRequest(route, http.MethodPost, "/project/Some%2FOther%2FProject/build/222")
	assert.Equal(t, http.StatusOK, w.Code)
}

// PerformRequest for testing gin router.
func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	// for _, h := range headers {
	// 	req.Header.Add(h.Key, h.Value)
	// }
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
