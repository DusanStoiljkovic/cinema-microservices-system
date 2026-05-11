# Cinema Microservices System

Cinema Microservices System is a backend project built with Go using a microservices-style architecture.

The system is designed for managing users, movies, genres, halls, projections, tickets, and cinema bookings through multiple independent services connected through an API Gateway.

---

<img width="1518" height="1036" alt="dbDiagram" src="https://github.com/user-attachments/assets/5cc750e0-5f61-426d-a25d-e746ba7d9a8e" />


## Architecture

The system is divided into several services:

| Service | Responsibility |
|---|---|
| API Gateway | Single entry point for client requests |
| User Service | User registration, login, authentication, user profile and user administration |
| Movie Service | Movies, genres, and movie-genre relations |
| Booking Service | Halls, projections, tickets, and orders |

Each service has its own internal structure with handlers, services, repositories, DTOs, models, and utilities.

---

## Tech Stack

- Go
- Chi Router
- GORM
- MySQL
- Docker
- Docker Compose
- JWT Authentication
- REST API
- Microservices-style architecture

---

## Project Structure

```text
cinema-microservices-system/
│
├── api-gateway/
│   └── cmd/
│
├── user-service/
│   ├── cmd/
│   └── internal/
│       ├── auth/
│       ├── dto/
│       ├── handlers/
│       ├── mapper/
│       ├── middleware/
│       ├── models/
│       ├── repository/
│       ├── routes/
│       ├── services/
│       └── utils/
│
├── movie-service/
│   ├── cmd/
│   └── internal/
│       ├── dto/
│       ├── handler/
│       ├── mapper/
│       ├── middleware/
│       ├── models/
│       ├── repository/
│       ├── routes/
│       ├── service/
│       └── utils/
│
├── booking-service/
│   ├── cmd/
│   └── internal/
│       ├── auth/
│       ├── dto/
│       ├── handler/
│       ├── mapper/
│       ├── middleware/
│       ├── models/
│       ├── repository/
│       ├── routes/
│       ├── service/
│       └── utils/
│
├── docker-compose.yaml
└── README.md
```

---

## Running the Project

Start the whole system with Docker Compose:

```bash
docker compose up --build
```

Stop all containers:

```bash
docker compose down
```

Remove containers and volumes:

```bash
docker compose down -v
```

---

## Environment Variables

Each service should have its own `.env` file.

Example files should be used as templates:

```text
api-gateway.env.example
user-service.env.example
movie-service.env.example
booking-service.env.example
```

Do not commit real `.env` files with secrets to GitHub.

---

## API Gateway

All requests should go through the API Gateway.

Base URL:

```text
http://localhost:8080
```

Main API route groups:

```text
/api/users
/api/movies
/api/genres
/api/halls
/api/projections
/api/tickets
/api/orders
```

The API Gateway removes the `/api` prefix before forwarding requests to the correct service.

Example:

```text
GET /api/movies
```

is forwarded internally to:

```text
GET /movies
```

---

# API Routes

## Authentication

Some routes require JWT authentication.

Login returns a JWT token:

```json
{
  "jwt": "your.jwt.token"
}
```

Use the token in protected routes:

```text
Authorization: Bearer your.jwt.token
```

---

# User Service

Base path:

```text
/api/users
```

## User Routes

| Method | Endpoint | Description | Auth |
|---|---|---|---|
| POST | `/api/users/register` | Register new user | Public |
| POST | `/api/users/login` | Login user | Public |
| GET | `/api/users/me` | Get current authenticated user | JWT |
| GET | `/api/users/{id}` | Get user by ID | JWT |
| PUT | `/api/users` | Update user | JWT |
| GET | `/api/users` | Get all users | JWT + Admin |
| DELETE | `/api/users/{id}` | Delete user by ID | JWT + Admin |

---

## Register User

```http
POST /api/users/register
```

Request body:

```json
{
  "name": "Dusan Stoiljkovic",
  "email": "dusan@gmail.com",
  "password": "dusan123"
}
```

---

## Login User

```http
POST /api/users/login
```

Request body:

```json
{
  "email": "dusan@gmail.com",
  "password": "dusan123"
}
```

Expected response:

```json
{
  "jwt": "your.jwt.token"
}
```

---

## Get Current User

```http
GET /api/users/me
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

## Get User By ID

```http
GET /api/users/{id}
```

Example:

```http
GET /api/users/1
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

## Get All Users

```http
GET /api/users
```

Headers:

```text
Authorization: Bearer admin.jwt.token
```

Required role:

```text
admin
```

---

## Update User

```http
PUT /api/users
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

Request body example:

```json
{
  "id": 1,
  "name": "Dusan Stoiljkovic",
  "email": "dusan@gmail.com",
  "password": "newpassword123",
  "role": "user"
}
```

> Note: The current route is `PUT /api/users`, not `PUT /api/users/{id}`.

---

## Delete User

```http
DELETE /api/users/{id}
```

Example:

```http
DELETE /api/users/1
```

Headers:

```text
Authorization: Bearer admin.jwt.token
```

Required role:

```text
admin
```

---

# Movie Service

Base paths:

```text
/api/movies
/api/genres
```

---

## Movie Routes

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/movies` | Get all movies / search movies |
| GET | `/api/movies/{id}` | Get movie by ID |
| GET | `/api/movies/{id}/genres` | Get genres for movie |
| POST | `/api/movies` | Create movie |
| POST | `/api/movies/{movie_id}/genres/{genre_id}` | Add genre to movie |
| PUT | `/api/movies/{id}` | Update movie |
| DELETE | `/api/movies/{id}` | Delete movie |
| DELETE | `/api/movies/{movie_id}/genres/{genre_id}` | Remove genre from movie |

---

## Get All Movies

```http
GET /api/movies
```

Optional query params:

| Query Param | Example | Description |
|---|---|---|
| `limit` | `20` | Number of movies to return |
| `offset` | `0` | Pagination offset |
| `sort` | `asc` / `desc` | Sort order |
| `genre` | `action` | Filter by genre name |
| `min_year` | `2020` | Minimum movie year |
| `max_year` | `2030` | Maximum movie year |
| `min_rating` | `7.0` | Minimum rating |

Example:

```http
GET /api/movies?genre=action&min_year=2020&max_year=2030
```

Another example:

```http
GET /api/movies?limit=10&offset=0&sort=desc&min_rating=7.5
```

---

## Get Movie By ID

```http
GET /api/movies/{id}
```

Example:

```http
GET /api/movies/1
```

---

## Get Movie Genres

```http
GET /api/movies/{id}/genres
```

Example:

```http
GET /api/movies/1/genres
```

---

## Create Movie

```http
POST /api/movies
```

Request body:

```json
{
  "title": "American Gangster",
  "description": "A thief who steals corporate secrets through dream-sharing technology.",
  "year": 2002,
  "image_url": "https://example.com/inception.jpg",
  "duration": 180,
  "rating": 7.0,
  "genre_ids": [1, 2, 6]
}
```

---

## Update Movie

```http
PUT /api/movies/{id}
```

Example:

```http
PUT /api/movies/1
```

Request body:

```json
{
  "title": "American Gangster",
  "description": "A thief who steals corporate secrets through dream-sharing technology.",
  "year": 2002,
  "image_url": "https://example.com/inception.jpg",
  "duration": 200,
  "rating": 8.0,
  "genre_ids": [1, 2]
}
```

---

## Delete Movie

```http
DELETE /api/movies/{id}
```

Example:

```http
DELETE /api/movies/3
```

---

## Add Genre To Movie

```http
POST /api/movies/{movie_id}/genres/{genre_id}
```

Example:

```http
POST /api/movies/1/genres/3
```

---

## Remove Genre From Movie

```http
DELETE /api/movies/{movie_id}/genres/{genre_id}
```

Example:

```http
DELETE /api/movies/1/genres/2
```

---

# Genres

Base path:

```text
/api/genres
```

## Genre Routes

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/genres` | Get all genres |
| GET | `/api/genres/{id}` | Get genre by filter |
| POST | `/api/genres` | Create genre |
| PUT | `/api/genres/{id}` | Update genre |
| DELETE | `/api/genres/{id}` | Delete genre |

---

## Get All Genres

```http
GET /api/genres
```

---

## Get Genre By Filter

```http
GET /api/genres/{id}
```

Example:

```http
GET /api/genres/1
```

> Note: The route exists as `/api/genres/{id}`, but the current handler reads a `GenreFilter` from the request body. If you want a classic `GET /api/genres/{id}`, the handler should be adjusted to use the URL param.

---

## Create Genre

```http
POST /api/genres
```

Request body:

```json
{
  "name": "romance"
}
```

---

## Update Genre

```http
PUT /api/genres/{id}
```

Example:

```http
PUT /api/genres/1
```

Request body:

```json
{
  "name": "drama"
}
```

---

## Delete Genre

```http
DELETE /api/genres/{id}
```

Example:

```http
DELETE /api/genres/2
```

---

# Booking Service

Base paths:

```text
/api/halls
/api/projections
/api/tickets
/api/orders
```

---

# Halls

## Hall Routes

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/halls` | Get all halls |
| GET | `/api/halls/{id}` | Get hall by ID |
| POST | `/api/halls` | Create hall |
| PUT | `/api/halls/{id}` | Update hall |
| DELETE | `/api/halls/{id}` | Delete hall |

---

## Get All Halls

```http
GET /api/halls
```

---

## Get Hall By ID

```http
GET /api/halls/{id}
```

Example:

```http
GET /api/halls/4
```

---

## Create Hall

```http
POST /api/halls
```

Request body:

```json
{
  "name": "Sala 3",
  "location": "Paracin",
  "capacity": 100
}
```

---

## Update Hall

```http
PUT /api/halls/{id}
```

Example:

```http
PUT /api/halls/4
```

Request body:

```json
{
  "name": "Sala 1",
  "location": "Paracin",
  "capacity": 500
}
```

---

## Delete Hall

```http
DELETE /api/halls/{id}
```

Example:

```http
DELETE /api/halls/1
```

---

# Projections

## Projection Routes

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/projections` | Get all projections |
| GET | `/api/projections/{id}` | Get projection by ID |
| GET | `/api/projections/movie/{movie_id}` | Get projections by movie ID |
| POST | `/api/projections` | Create projection |
| PUT | `/api/projections/{id}` | Update projection |
| DELETE | `/api/projections/{id}` | Delete projection |

---

## Get All Projections

```http
GET /api/projections
```

---

## Get Projection By ID

```http
GET /api/projections/{id}
```

Example:

```http
GET /api/projections/1
```

---

## Get Projections By Movie ID

```http
GET /api/projections/movie/{movie_id}
```

Example:

```http
GET /api/projections/movie/1
```

---

## Create Projection

```http
POST /api/projections
```

Request body:

```json
{
  "movie_id": 1,
  "hall_id": 4,
  "start_time": "2026-05-05T18:00:00+02:00",
  "end_time": "2026-05-05T20:00:00+02:00",
  "price": 500
}
```

---

## Update Projection

```http
PUT /api/projections/{id}
```

Example:

```http
PUT /api/projections/3
```

Request body:

```json
{
  "movie_id": 1,
  "hall_id": 4,
  "start_time": "2026-05-05T18:00:00+02:00",
  "end_time": "2026-05-05T20:00:00+02:00",
  "price": 600
}
```

---

## Delete Projection

```http
DELETE /api/projections/{id}
```

Example:

```http
DELETE /api/projections/3
```

---

# Tickets

## Ticket Routes

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/tickets` | Get all tickets |
| GET | `/api/tickets/{id}` | Get ticket by ID |
| GET | `/api/tickets/users/{id}` | Get tickets by user ID |
| GET | `/api/tickets/projections/{id}` | Get tickets by projection ID |
| POST | `/api/tickets` | Create ticket |
| DELETE | `/api/tickets/{id}` | Delete ticket |

---

## Get All Tickets

```http
GET /api/tickets
```

---

## Get Ticket By ID

```http
GET /api/tickets/{id}
```

Example:

```http
GET /api/tickets/1
```

---

## Get Tickets By User ID

```http
GET /api/tickets/users/{id}
```

Example:

```http
GET /api/tickets/users/1
```

---

## Get Tickets By Projection ID

```http
GET /api/tickets/projections/{id}
```

Example:

```http
GET /api/tickets/projections/1
```

---

## Create Ticket

```http
POST /api/tickets
```

Request body:

```json
{
  "projection_id": 1,
  "seat_number": 1
}
```

---

## Delete Ticket

```http
DELETE /api/tickets/{id}
```

Example:

```http
DELETE /api/tickets/1
```

---

# Orders

Base path:

```text
/api/orders
```

All order routes require JWT authentication.

Headers:

```text
Authorization: Bearer your.jwt.token
```

## Order Routes

| Method | Endpoint | Description | Auth |
|---|---|---|---|
| GET | `/api/orders` | Get all orders | JWT |
| GET | `/api/orders/{id}` | Get order by ID | JWT |
| GET | `/api/orders/users/{id}` | Get orders by user ID | JWT |
| GET | `/api/orders/me` | Get current user's orders | JWT |
| POST | `/api/orders` | Create order | JWT |
| PATCH | `/api/orders/pay/{id}` | Pay order | JWT |
| PATCH | `/api/orders/cancel/{id}` | Cancel order | JWT |
| DELETE | `/api/orders/{id}` | Delete order | JWT |

---

## Get All Orders

```http
GET /api/orders
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

## Get Order By ID

```http
GET /api/orders/{id}
```

Example:

```http
GET /api/orders/5
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

## Get Orders By User ID

```http
GET /api/orders/users/{id}
```

Example:

```http
GET /api/orders/users/1
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

## Get My Orders

```http
GET /api/orders/me
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

## Create Order

```http
POST /api/orders
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

Request body:

```json
{
  "tickets": [
    {
      "projection_id": 1,
      "seat_number": 1
    },
    {
      "projection_id": 1,
      "seat_number": 2
    }
  ]
}
```

---

## Pay Order

```http
PATCH /api/orders/pay/{id}
```

Example:

```http
PATCH /api/orders/pay/5
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

## Cancel Order

```http
PATCH /api/orders/cancel/{id}
```

Example:

```http
PATCH /api/orders/cancel/5
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

## Delete Order

```http
DELETE /api/orders/{id}
```

Example:

```http
DELETE /api/orders/5
```

Headers:

```text
Authorization: Bearer your.jwt.token
```

---

# Error Response Format

The application uses a centralized SafeError-style error handler.

Example error response:

```json
{
  "code": "INVALID_INPUT",
  "error": "Invalid request body"
}
```

Common error codes:

| Code | HTTP Status | Description |
|---|---:|---|
| `INVALID_INPUT` | 400 | Invalid request body, params, or validation |
| `AUTH_FAILED` | 401 | Authentication failed |
| `FORBIDDEN` | 403 | User does not have permission |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource conflict |
| `INTERNAL_ERROR` | 500 | Internal server error |

---

# Current Features

- User registration
- User login
- JWT authentication
- Role-based admin protection for selected user routes
- Get current user
- Get user by ID
- Get all users as admin
- Update user
- Delete user as admin
- Movie CRUD
- Movie search and filtering through query params
- Genre CRUD
- Movie-genre relations
- Hall CRUD
- Projection CRUD
- Ticket creation, ticket queries, and ticket deletion
- Order creation
- Get all orders
- Get order by ID
- Get orders by user ID
- Get current user's orders
- Pay order
- Cancel order
- Delete order
- API Gateway routing
- MySQL databases
- Docker Compose setup
- SafeError-style centralized error handling

---

# Implementation Notes

## User Update Route

The current user update route is:

```http
PUT /api/users
```

It is not currently:

```http
PUT /api/users/{id}
```

If you want a more REST-style endpoint, change the route to:

```go
r.Put("/{id}", middleware.ErrorHandler(userHandler.HandleUpdateUser))
```

and then read the user ID from the URL param inside the handler.

---

## Genre By ID Route

The current route exists as:

```http
GET /api/genres/{id}
```

However, the current handler reads a `GenreFilter` from the request body instead of using the `{id}` URL param.

A cleaner implementation would be to parse the `id` from the URL param and call the service with:

```go
&dto.GenreFilter{ID: &id}
```

---

## API Gateway Booking Prefix

The API Gateway currently contains a `/api/bookings` proxy prefix, but the booking service does not currently define a `/bookings` route group.

The active booking-related route groups are:

```text
/api/halls
/api/projections
/api/tickets
/api/orders
```

---

# Future Improvements

- Add refresh tokens
- Add real logout with refresh-token invalidation or JWT blacklist
- Improve order authorization rules
- Add admin protection for booking-management routes if needed
- Add seat availability endpoints
- Add recommendation-service endpoints
- Add unit tests for service layer
- Add integration tests with test database
- Add database migrations
- Add health check endpoints
- Add CI/CD pipeline
- Add Swagger/OpenAPI documentation
- Improve logging with request IDs
- Add event-driven communication between services
