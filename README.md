# marumaru-backend
Backend for a universal marketplace

## Setup

```
go get -v -t -d ./...
go run api.go
```

The server starts up at `localhost:8080`

## Features

Plan to implement all short goals within deadline, and as many long term goals as possible. Please see [Vision doc](https://docs.google.com/document/d/1bldfOAaVAh2pxbKF_G5jPiawzW5ktfYV5UtP09rgYl0/edit#) for full list of features.


### Login and Register
----
**Signin**
  For user sign in.

* **URL**

  /user/login
  Make sure set {"Cookie": #cookie_token#} in your header.

* **Method:**
  
  `POST`
  
* **Success Response:**
  
  * **Code:** 200 <br />

* **Error Response:**

  * **Code:** 401 <br />
  * **Code:** 500 <br />

----
**SignUp**
  For user sign up.

* **URL**

  /user/signup

* **Method:**
  
  `POST`

* **Data Params**

  A JSON of a User struct,
  where username, password, email must be nonempty.
* **Success Response:**
  
  * **Code:** 200 <br />
  In response header, {"Set-Cookie": #cookie_token#} included.

* **Error Response:**

  * **Code:** 401 <br />
  * **Code:** 400 <br />
  * **Code:** 500 <br />

----

**Refresh**
  Refresh the token in the background by the client application.
* **URL**

  /user/refresh
  Make sure set {"Cookie": #cookie_token#} in your header.

* **Method:**
  
  `POST`

* **Success Response:**
  
  * **Code:** 200 <br />
  In response header, A token {"Set-Cookie": #cookie_token#} included.

* **Error Response:**

  * **Code:** 401 <br />
  * **Code:** 400 <br />
  * **Code:** 500 <br />

----

**GetUserByCookie**
GetUserByCookie returns the whole user struct by your cookie
* **URL**

  /user
  Make sure set {"Cookie": #cookie_token#} in your header.

* **Method:**
  
  `GET`

* **Success Response:**
  
  * **Code:** 200 <br />
  A JSON file of User struct included in body.
* **Error Response:**

  * **Code:** 401 <br />
  * **Code:** 400 <br />
  * **Code:** 500 <br />
----

### Get products, or get single product

**getAllProductsHandler**
getAllProductsHandler returns all the products in the database

- **URL**

  /products

- **Method:**

  `GET`

- **Success Response:**

  - **Code:** 200 <br />
    A JSON file of Product list.

- **Error Response:**

  - **Code:** 500 <br />

----

**getProductByIDHandler**

getProductByIDHandler returns the product with the given product id.

* **URL**

  /product/:id

  id is the product id, it should be an uuid.

* **Method:**

  `GET`

* **Success Response:**

  - **Code:** 200 <br />
    A JSON file of the Product.

* **Error Response:**

  - **Code:** 500 <br />

### Add new product

**getProductByIDHandler**

getProductByIDHandler returns the product with the given product id.

- **URL**

  /products

  Make sure set {"Cookie": #cookie_token#} in your header.

- **Method:**

  `POST`

* **Data Params**

  A JSON of a Product struct,
  where Photos , Name, Description,  Price, TimeDuration, Tags should be included.

- **Success Response:**

  - **Code:** 200 <br />
  
- **Error Response:**

  - **Code:** 500 <br />

### Buy product

### Recommend products

### Search products

### Edit listings

### User profile for history

**getProductByUserIDHandler**

**getOrderByUserIDHandler**

**getOrderByIDHandler**

**getOrderByProductIDHandler**

### Reviews / comments

### Filter/ Sort (date / price / etc)


### Automatic tagging