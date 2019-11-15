package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/merumaru/marumaru-backend/models"
	. "github.com/merumaru/marumaru-backend/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRequestWithBody(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRequestWithBodyAndHeader(r http.Handler, method, path string, body io.Reader, key string, value string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set(key, value)
	w := httptest.NewRecorder()
	// w.Header().Set(key, value)
	r.ServeHTTP(w, req)
	return w
}

var _ = Describe("Server", func() {
	var (
		router   *gin.Engine
		response *httptest.ResponseRecorder
		client   *mongo.Client
	)

	BeforeEach(func() {
		router = CreateRouter()
		// StartServer(router)
		client = Connect2DB()
	})

	Describe("login API", func() {
		var (
			collection *mongo.Collection
		)
		BeforeEach(func() {
			collection = client.Database("testing").Collection("users")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			// drop the table
			collection.Drop(ctx)
		})
		Describe("register API at /signup", func() {
			BeforeEach(func() {
				user := User{Username: "user1", Password: "password1", Email: "www"}
				ujson, _ := json.Marshal(user)
				body := bytes.NewReader(ujson)
				response = performRequestWithBody(router, "POST", "/signup", body)
			})

			It("Add a user", func() {
				filter := bson.M{"username": "user1"}
				result := User{}
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				err := collection.FindOne(ctx, filter).Decode(&result)
				Expect(response.Code).To(Equal(200))
				Expect(err).To(BeNil())
				Expect(result.Password).To(Equal("password1"))
				Expect(result.Email).To(Equal("www"))
			})
		})

		Describe("register API at /signup", func() {
			BeforeEach(func() {
				user := User{Username: "user1", Password: "password1", Email: ""}
				ujson, _ := json.Marshal(user)
				body := bytes.NewReader(ujson)
				response = performRequestWithBody(router, "POST", "/signup", body)
			})

			It("Add a invalid user", func() {
				Expect(response.Code).To(Equal(400))
			})
		})

		Describe("register API at /signup", func() {
			BeforeEach(func() {
				user := User{Username: "user1", Password: "password1", Email: ""}
				ujson, _ := json.Marshal(user)
				body := bytes.NewReader(ujson)
				_ = performRequestWithBody(router, "POST", "/signup", body)
				response = performRequestWithBody(router, "POST", "/signup", body)
			})

			It("Refuse add the same user twice", func() {
				Expect(response.Code).To(Equal(400))
			})
		})

		Describe("Signin API at /login", func() {
			BeforeEach(func() {
				// sign up a user
				user := User{Username: "user1", Password: "password1", Email: "www"}
				ujson, _ := json.Marshal(user)
				body := bytes.NewReader(ujson)
				_ = performRequestWithBody(router, "POST", "/signup", body)
				// sign in
				cred := Credentials{Password: "password1", Username: "user1"}
				cjson, _ := json.Marshal(cred)
				body = bytes.NewReader(cjson)
				response = performRequestWithBody(router, "POST", "/login", body)
				println(string(response.Body.Bytes()))
			})

			It("Ok!", func() {
				Expect(response.Code).To(Equal(200))
			})
		})

		Describe("Signin API at /login with wrong pwd", func() {
			BeforeEach(func() {
				// sign up a user
				user := User{Username: "user1", Password: "password11", Email: "www"}
				ujson, _ := json.Marshal(user)
				body := bytes.NewReader(ujson)
				_ = performRequestWithBody(router, "POST", "/signup", body)
				// sign in
				cred := Credentials{Password: "password1", Username: "user1"}
				cjson, _ := json.Marshal(cred)
				body = bytes.NewReader(cjson)
				response = performRequestWithBody(router, "POST", "/login", body)
				println(string(response.Body.Bytes()))
			})

			It("No!", func() {
				Expect(response.Code).To(Equal(401))
			})
		})

		Describe("Signin API at /login with no user", func() {
			BeforeEach(func() {
				// sign in
				cred := Credentials{Password: "password1", Username: "user1"}
				cjson, _ := json.Marshal(cred)
				body := bytes.NewReader(cjson)
				response = performRequestWithBody(router, "POST", "/login", body)
				println(string(response.Body.Bytes()))
			})

			It("No such user!", func() {
				Expect(response.Code).To(Equal(404))
			})
		})

		Describe("Welcome API at /welcome", func() {
			BeforeEach(func() {
				// sign up a user
				user := User{Username: "user1", Password: "password1", Email: "www"}
				ujson, _ := json.Marshal(user)
				body := bytes.NewReader(ujson)
				_ = performRequestWithBody(router, "POST", "/signup", body)
				// sign in
				cred := Credentials{Password: "password1", Username: "user1"}
				cjson, _ := json.Marshal(cred)
				body = bytes.NewReader(cjson)
				response = performRequestWithBody(router, "POST", "/login", body)
				// welcome
				response = performRequestWithBodyAndHeader(router, "GET", "/welcome", nil,
					"Cookie", response.Header().Get("Set-Cookie"))
				println(string(response.Body.Bytes()))
			})

			It("Ok!", func() {
				Expect(response.Code).To(Equal(200))
			})
		})

		Describe("Welcome API at /welcome without cookie", func() {
			BeforeEach(func() {
				// welcome
				response = performRequestWithBodyAndHeader(router, "GET", "/welcome", nil,
					"Cookie", response.Header().Get("Set-Cookie"))
				println(string(response.Body.Bytes()))
			})

			It("No cookie!", func() {
				Expect(response.Code).To(Equal(401))
			})
		})

		Describe("Refresh API at /welcome", func() {
			BeforeEach(func() {
				// sign up a user
				user := User{Username: "user1", Password: "password1", Email: "www"}
				ujson, _ := json.Marshal(user)
				body := bytes.NewReader(ujson)
				_ = performRequestWithBody(router, "POST", "/signup", body)
				// sign in
				cred := Credentials{Password: "password1", Username: "user1"}
				cjson, _ := json.Marshal(cred)
				body = bytes.NewReader(cjson)
				response = performRequestWithBody(router, "POST", "/login", body)
				// welcome
				response = performRequestWithBodyAndHeader(router, "POST", "/refresh", nil,
					"Cookie", response.Header().Get("Set-Cookie"))
				println(string(response.Body.Bytes()))
			})

			It("A new token will only be issued if the old token is within 30 seconds of expiry!", func() {
				Expect(response.Code).To(Equal(400))
			})
		})

		Describe("Get user API at /usercookie", func() {
			var user User
			BeforeEach(func() {
				// sign up a user
				user = User{Username: "user1", Password: "password1", Email: "www"}
				ujson, _ := json.Marshal(user)
				body := bytes.NewReader(ujson)
				_ = performRequestWithBody(router, "POST", "/signup", body)
				// sign in
				cred := Credentials{Password: "password1", Username: "user1"}
				cjson, _ := json.Marshal(cred)
				body = bytes.NewReader(cjson)
				response = performRequestWithBody(router, "POST", "/login", body)
				// welcome
				response = performRequestWithBodyAndHeader(router, "GET", "/usercookie", nil,
					"Cookie", response.Header().Get("Set-Cookie"))
			})

			It("Ok!", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Return the user", func() {
				var actualUser User
				json.Unmarshal(response.Body.Bytes(), &actualUser)
				Expect(actualUser.Username).To(Equal(user.Username))
			})
		})
	})

	// Describe("product API", func() {
	// 	var (
	// 		collection *mongo.Collection
	// 	)
	// 	BeforeEach(func() {
	// 		collection = client.Database("testing").Collection("users")
	// 		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// 		// drop the table
	// 		collection.Drop(ctx)
	// 	})
	// 	Describe("addProductHandler API", func ()  {
	// 		BeforeEach(func() {

	// 		})
	// 	})
	// })
})
