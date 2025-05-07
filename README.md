# Go REST API Project (Social Media Case)

A scalable and modular REST API built with Go using a clean architecture pattern (Handler → Service → Repository). This application demonstrates CRUD functionality and follows best practices for structuring Go web applications.

- [🧠 Application Overview](#-application-overview)
- [🚀 Features](#-features)
- [🛠️ Tech Stack](#-tech-stack)
- [🧰 How to Run the Application](#-how-to-run-the-application)
- [📌 API Documentation](#-api-documentation)
- [🔐 Authentication](#-authentication)
- [🔐 Authorization](#-authorization)
- [🧪 Testing](#-testing)
- [🗂️ Project Structure](#-project-structure)

---

## 🧠 Application Overview

This is a RESTful API server built in Go Programming Language, following clean architecture principles. Each layer has a clear responsibility:

- **Handler**: receives HTTP requests and sends responses
- **Service**: contains business logic
- **Repository**: handles communication with the database

This structure improves scalability, testability, and maintainability.

---

## 🚀 Features

- RESTful API with structured endpoints
- Modular architecture with separation of concerns (Handler → Service → Repository)
- Middleware support (e.g., Basic Authentication, Bearer Authentication, Permission, and Rate Limiting)
- Caching using Redis (In-Memory database)
- PostgreSQL as relational database
- GORM as ORM library for database interaction
- Database table relations (One-to-one, One-to-many, Many-to-many)
- Testing support
- Dockerized setup for easy deployment
- Manual / local setup without docker

---

## 🛠️ Tech Stack

This project uses the following technologies:

- **Go** (`https://github.com/golang`) — Main programming language.
- **PostgreSQL** (`https://github.com/postgres/postgres`) — Relational database.
- **Redis** (`https://github.com/redis/redis`) — In-Memory database.
- **Chi Router** (`https://github.com/go-chi/chi`) — Lightweight and idiomatic HTTP router.
- **GORM** (`https://github.com/go-gorm/gorm`) — ORM library for database interaction.
- **PGX** (`https://github.com/jackc/pgx`) — PostgreSQL driver for GORM.
- **Google UUID** (`https://github.com/google/uuid`) — For generating UUIDs.
- **Viper** (`https://github.com/spf13/viper`) — Loads environment variables from a `.env` file.
- **Testify** (`https://github.com/stretchr/testify`) — Testing framework.
- **Go-sqlmock** (`https://github.com/DATA-DOG/go-sqlmock`) — Mocking database query.
- **Go Validator** (`https://github.com/go-playground/validator`) — Validating body request.
- **Go JWT** (`https://github.com/golang-jwt/jwt`) — Generate & parse JSON Web Token.
- **Logrus** (`https://github.com/sirupsen/logrus`) — Logging.
- **Architecture** — Clean architecture: Handler → Service → Repository.
- **Docker & Docker Compose** — For containerized development.

---

## 🧰 How to Run the Application

There are two ways to run this application: using Docker (recommended for easy setup) or manually running it on your local machine. Choose one based on your preference.

#### 📦 Option 1: Run with Docker (Recommended)

> This is the easiest way to get the app running. It uses Docker and Docker Compose to set up everything automatically.

**1. Ensure Docker and Docker Compose are installed** on your machine.
**2. Clone the repository**:

```bash
$ git clone https://github.com/aryaadinulfadlan/go-social-api.git
$ cd go-social-api
```

**3. Make sure you are on the `main` branch. Build the application:**

```bash
$ docker-compose up --build
```

**4. Once the app is successfully built, run the database migration:**

```bash
$ docker-compose run --rm migrate
```

**5. You are all set. Feel free to access any of the available endpoints using Postman or curl: http://localhost:4000.**

#### 📦 Option 2: Run Locally (Without Docker)

> If you prefer running the application manually without Docker, follow these steps. Ensure you have **Go v1.24.1**, **PostgreSQL v16.4**, and **Redis v7.2.8** installed on your machine.

**1. Clone the repository**:

```bash
$ git clone https://github.com/aryaadinulfadlan/go-social-api.git
$ cd go-social-api
```

**2. Checkout to `local` branch:**

```bash
$ git checkout local
```

**3. Make sure PostgreSQL is running on your local machine.**
**4. Start Redis on your local machine:**

```bash
$ redis-server
```

**5. On another terminal tab, you can test your Redis connection with the following command:**

```bash
$ redis-cli
$ ping
```

If it returns PONG, it means you're successfully connected to Redis.
**6. Adjust the database URL in the `.env `file to match your own**
**7. Create a new database on your machine**
**8. I am using `golang-migrate`, you can choose your own. Run the database migration:**

```bash
$ make migrate_up
```

**9. Run the application:**

```bash
$ go run ./cmd/server/
```

**10. You are all set. Feel free to access any of the available endpoints using Postman or curl: http://localhost:4000**

---

## 📌 API Documentation

This section outlines the available API endpoints, including request parameters, request body, authentication requirements, and response formats.

#### 🔓 Unauthenticated Endpoints

##### 1. `GET /v1/ping` - Just for checking

##### 2. `POST /v1/auth/sign-up` - Register to the application

**Request Body:**

```json
{
  "name": "John Doe",
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

##### 3. `POST /v1/auth/sign-in` — Login to the application

**Request Body:**

```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

##### 4. `POST /v1/auth/resend-activation` — Resend / Generate token activation if previously token is expired

**Request Body:**

```json
{
  "email": "john@example.com"
}
```

##### 5.`PUT /v1/auth/activate/{token}` — Activate a user using token activation (**default user is inactive**)

**Request Params: token (string)**

#### 🔐 Authenticated Endpoints

##### 1. `GET /v1/basic` — Trying Basic Authenctication scheme

##### 2. `GET /v1/bearer` — Trying Bearer Authenctication scheme

##### 3. `POST /v1/posts` - Create a new post

**Request Body:**

```json
{
  "title": "Post Title",
  "content": "Post Content",
  "tags": ["haha", "hihi"]
}
```

##### 4. `GET /v1/posts/{postId}` — Get detail of a specific post

**Request Params: postId (uuid)**

##### 5. `PATCH /v1/posts/{postId}` — Update a specific post

**Request Params: postId (uuid)**
**Request Body:**

```json
{
  "title": "Post Updated",
  "content": "Post Content Updated",
  "tags": ["huhu", "hihi"]
}
```

##### 6. `DELETE /v1/posts/{postId}` — Delete a specific post

**Request Params: postId (uuid)**

##### 7. `GET /v1/users/{userId}` — Get detail of a specific user

**Request Params: userId (uuid)**

##### 8. `POST /v1/users/{userId}/follow` — Follow or Unfollow a specific user

**Request Params: userId (uuid)**

##### 9. `GET /v1/users/{userId}/followers` — Get the list of followers of a specific user

**Request Params: userId (uuid)**

##### 10. `GET /v1/users/{userId}/following` — Get the list of accounts followed by a specific user

**Request Params: userId (uuid)**

##### 11. `GET /v1/users/feed` — Get a user's feed posts (a combination of their own posts and posts from users they follow)

**Query String**

```json
{
  "per_page": 10,
  "page": 1,
  "sort": "DESC",
  "tags": ["tag-one", "tag-two"],
  "search": "keyword",
  "since": "2025-04-10",
  "until": "2025-04-11"
}
```

##### 12. `POST /v1/comments/{postId}` — Add a comment to a post

**Request Params: postId (uuid)**
**Request Body:**

```json
{
  "content": "I commented a post"
}
```

---

## 🔐 Authentication

This application supports two authentication methods: **Basic Authentication** and **Bearer Token Authentication (JWT)**.

#### 🔑 Basic Authentication

You can authenticate by sending your `username` and `password` using the `Authorization` header with the **Basic** scheme. Alternatively, for convenience, I've provided the username, password, and the pre-encoded Base64 string.

- username: **arya** (in the .env file)
- password: **arya123** (in the .env file)
- base64: **YXJ5YTphcnlhMTIz**

I only set up one endpoint that requires Basic Authentication, which is `/v1/basic`. You can try accessing this endpoint.

#### 🔑 Bearer Authentication

Once logged in using valid credentials, the server will return a JWT token. Use this token for subsequent requests to access protected endpoints.

Most endpoints are protected using Bearer Authentication. You can find the list in the **API Documentation** section.

---

## 🔐 Authorization

This application implements **role-based access control** with two user types:

- **Admin**
- **Regular User**

Each user type has specific access rights to certain API endpoints.
Here are the accounts you can use to log in to this application:

- **Princess Diana (Admin User):**
  - Email: `princess_diana@gmail.com`
  - Password: `diana123`
- **Clark Kent (Regular User):**
  - Email: `clark_kent@gmail.com`
  - Password: `clark123`
- **Bruce Wayne (Regular User):**
  - Email: `bruce_wayne@gmail.com`
  - Password: `bruce123`

#### 👤 Regular User

Regular users have **limited access** to the system. They are allowed to access only specific endpoints.

##### ✅ Allowed Endpoints for Regular Users:

- `GET /v1/ping` — Just for checking
- `GET /v1/basic` — Trying Basic Authenctication scheme
- `GET /v1/bearer` — Trying Bearer Authenctication scheme
- `POST /v1/posts` — Create a new post
- `GET /v1/posts/{postId}` — Get detail of a specific post
- `PATCH /v1/posts/{postId}` — Update a specific post (**only its own post**)
- `DELETE /v1/posts/{postId}` — Delete a specific post (**only its own post**)
- `GET /v1/users/{userId}` — Get detail of a specific user
- `POST /v1/users/{userId}/follow` — Follow or Unfollow a specific user
- `GET /v1/users/{userId}/followers` — Get the list of followers of a specific user
- `GET /v1/users/{userId}/following` — Get the list of accounts followed by a specific user
- `GET /v1/users/feed` — Get a user's feed posts (a combination of their own posts and posts from users they follow)
- `POST /v1/comments/{postId}` — Add a comment to a post

Any attempt to access restricted endpoints will return a **403 Forbidden** response.

#### 👑 Admin

Admin users have **full access** to all endpoints, including all regular user endpoints plus the following additional ones:

##### ✅ Allowed Endpoints fo Admin Users:

- All regular users endpoints
- `PATCH /v1/posts/{postId}` — Someone else's post
- `DELETE /v1/posts/{postId}` — Someone else's post
- `DELETE /v1/users/{userId}` — Delete specific user

#### 🔒 How Authorization Works

Once authenticated (via JWT Authentication), the server will check the user's role / permission and determine whether the user is allowed to access a specific route.

If the user lacks permission, the server will respond with:

```json
{
  "code": 403,
  "status": "FORBIDDEN",
  "message": "you do not have permission to access this resource."
}
```

##### ✅ Allowed Endpoints for Unauthenticated Users:

- `POST /v1/auth/sign-up` — Register to the application
- `POST /v1/auth/sign-in` — Login to the application
- `POST /v1/auth/resend-activation` — Resend / Generate token activation if previously token is expired
- `PUT /v1/auth/activate/{token}` — Activate a user using token activation (**default user is inactive**)

---

## 🧪 Testing

This project supports two types of testing to ensure the correctness and stability of the application:

**1. Integration Testing** (using real database and redis)
**2. Mock Query / Unit Testing** (using `go-sqlmock`)

#### 🔁 Integration Testing

Integration tests are used to test the behavior of the system with real dependencies such as **PostgreSQL** and **Redis**.

#### ✅ Characteristics:

- Requires a real PostgreSQL database instance
- Requires a running Redis server
- Tests the full flow of requests including middleware, handler, service, and repository layers

#### 🔁 Mock Query / Unit Testing

Unit tests use go-sqlmock to simulate database behavior without connecting to a real database.

#### ✅ Characteristics:

- Fast and lightweight
- To verify SQL queries and expected results.
- Tests specific units like repository (database query)

#### 🔸 How to Run:

You can run each test one by one, or run all tests at once:

- Navigate to the /test folder

```json
   $ cd test
```

- Run all of test cases

```json
   $ go test -v
```

---

## 🗂️ Project Structure

This project follows a **modular and layered architecture** designed for scalability, maintainability, and testability. The key layers are:

- **Handler Layer**: Responsible for HTTP request and response handling (controllers).
- **Service Layer**: Contains business logic.
- **Repository Layer**: Handles database access and queries.
- **Middleware**: Encapsulates cross-cutting concerns like authentication, authorization, and logging.
- **Models & DTOs**: Define data structures used across layers.
- **Tests**: Contains both unit and integration tests.
