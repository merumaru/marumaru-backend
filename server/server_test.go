package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/merumaru/marumaru-backend/data"
	"github.com/merumaru/marumaru-backend/models"
	. "github.com/merumaru/marumaru-backend/models"
	. "github.com/merumaru/marumaru-backend/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func insertTestProduct(client *mongo.Client) primitive.ObjectID {
	product := Product{Name: "Test Product"}
	product.Start, _ = time.Parse("20060102", "20000101")
	product.End, _ = time.Parse("20060102", "20200202")
	product.Photos = make([]string, 1)
	product.Tags = make([]string, 1)
	product.ID = primitive.NewObjectID()
	data.AddProduct(client, &product)
	return product.ID
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

	Describe("product API", func() {
		var (
			collection *mongo.Collection
			id         primitive.ObjectID
			cookie     string
		)
		BeforeEach(func() {
			collection = client.Database("testing").Collection("users")
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			// drop the table
			collection.Drop(ctx)

			client.Database("testing").Collection("products").Drop(nil)
			client.Database("testing").Collection("orders").Drop(nil)
			id = insertTestProduct(client)

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

			cookie = response.Header().Get("Set-Cookie")
		})

		Describe("List products API at /products", func() {
			BeforeEach(func() {
				response = performRequest(router, "GET", "/products")
			})

			It("Ok!", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Return the list of products", func() {
				actualProducts := make([]Product, 0)
				json.Unmarshal(response.Body.Bytes(), &actualProducts)
				log.Println(actualProducts)
				Expect(len(actualProducts)).To(Equal(1))
				Expect(actualProducts[0].Name).To(Equal("Test Product"))
			})
		})

		Describe("Add product API at /products", func() {
			BeforeEach(func() {
				product := Product{Name: "Test Product 2"}
				pjson, _ := json.Marshal(product)
				body := bytes.NewReader(pjson)
				response = performRequestWithBody(router, "POST", "/products", body)
			})

			It("Created!", func() {
				Expect(response.Code).To(Equal(201))
			})

			It("Return the list of products", func() {
				actualProducts, _ := data.GetAllProducts(client)
				log.Println(actualProducts)
				Expect(len(*actualProducts)).To(Equal(2))
				Expect((*actualProducts)[1].Name).To(Equal("Test Product 2"))
			})
		})

		Describe("Get product by ID API at /products/:id", func() {
			BeforeEach(func() {
				response = performRequest(router, "GET", "/products/"+id.Hex())
			})

			It("Ok!", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Return the product", func() {
				actualProduct := models.Product{}
				json.Unmarshal(response.Body.Bytes(), &actualProduct)
				Expect(actualProduct.Name).To(Equal("Test Product"))
			})
		})

		Describe("Edit product by ID API at /products/:id/edit", func() {
			BeforeEach(func() {
				product := Product{Description: "Test Description"}
				pjson, _ := json.Marshal(product)
				body := bytes.NewReader(pjson)
				response = performRequestWithBody(router, "PATCH", "/products/"+id.Hex()+"/edit", body)
			})

			It("Ok!", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Modify the product", func() {
				actualProduct, _ := data.GetProductByID(client, id.Hex())
				Expect(actualProduct.Description).To(Equal("Test Description"))
			})
		})

		Describe("Rent product by ID API at /products/:id/rent", func() {
			BeforeEach(func() {
				response = performRequestWithBodyAndHeader(router, "POST", "/products/"+id.Hex()+"/rent?startDate=2019-11-01&endDate=2019-12-02", nil, "Cookie", cookie)
			})

			It("Ok!", func() {
				Expect(response.Code).To(Equal(200))
			})

			It("Generate an order", func() {
				actualOrders, _ := data.GetOrderByProductID(client, id.Hex())
				log.Println(actualOrders)
				log.Println(id.Hex())
				Expect(len(*actualOrders)).To(Equal(1))
				expected, _ := time.Parse("20060102", "20191101")
				Expect((*actualOrders)[0].Start).To(Equal(expected))
			})
		})
	})
})
