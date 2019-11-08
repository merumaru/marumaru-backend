package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/merumaru/marumaru-backend/data"
	"github.com/merumaru/marumaru-backend/models"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var _ = Describe("Server", func() {
	var (
		router   *gin.Engine
		response *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		router = CreateRouter()
		// Since we modify lists in memory, we need to restore them to a clean state before every test
		data.Reload()
	})

	Describe("Version 1 API at /api/v1", func() {
		Describe("The / endpoint", func() {
			BeforeEach(func() {
				response = performRequest(router, "GET", "/api/v1/")
			})

			It("Returns with Status 200", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Returns the String 'Hello World'", func() {
				Expect(response.Body.String()).To(Equal("Hello World"))
			})
		})

		Describe("GET the /pokemon endpoint", func() {
			BeforeEach(func() {
				response = performRequest(router, "GET", "/api/v1/pokemon")
			})

			It("Returns with Status 200", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Returns all Pokemon", func() {
				// To easily test JSON responses, you should unmarshal the response into an actual type, then compare to what you expect
				var actual models.PokemonList
				json.Unmarshal(response.Body.Bytes(), &actual)
				Expect(actual).To(Equal(*data.Pokemon()))
			})
		})
	})
})
