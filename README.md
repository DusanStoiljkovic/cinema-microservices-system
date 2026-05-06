# Cinema Microservices System

Cinema Microservices System is a backend project built with Go using a microservices-style architecture.

The system is designed for managing users, movies, genres, halls, projections, tickets, and cinema bookings through multiple independent services connected through an API Gateway.

This README is aligned with the current Postman collection for the project.

---

## Diagram 
<img width="1518" height="1036" alt="dbDiagram" src="https://github.com/user-attachments/assets/cf5c6eff-8835-48b9-a26b-3428c8b7a0c4" />


## Architecture

The system is divided into several services:

| Service | Responsibility |
|---|---|
| API Gateway | Single entry point for client requests |
| User Service | User registration, login, authentication, and user profile |
| Movie Service | Movies, genres, and movie-genre relations |
| Booking Service | Halls, projections, tickets, seats, and orders |
| Recommendation Service | Movie recommendation endpoints - currently in progress |

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

```txt
cinema-microservices-system/
│
├── api-gateway/
│   └── cmd/
│
├── user-service/
│   ├── cmd/
│   └── internal/
│       ├── dto/
│       ├── handlers/
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
│       ├── models/
│       ├── repository/
│       ├── routes/
│       ├── service/
│       └── utils/
│
├── booking-service/
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
├── recommendation-service/
│   └── cmd/
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

```txt
api-gateway.env.example
user-service.env.example
movie-service.env.example
booking-service.env.example
recommendation-service.env.example
```

Do not commit real `.env` files with secrets to GitHub.

---

## API Gateway

All requests should go through the API Gateway.

Base URL:

```txt
http://localhost:8080
```

Main route groups:

```txt
/api/users
/api/movies
/api/genres
/api/halls
/api/projections
/api/tickets
/api/orders
/api/recommendations
```

---

## API Status Legend

| Status | Meaning |
|---|---|
| Completed | Request has a defined URL in the Postman collection |
| In progress | Request exists in Postman, but the URL is missing or empty |

---

# API Routes

## User Service

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | TBD | Get all users | In progress |
| GET | `/api/users/{user_id}` | Get user by ID | Completed |
| GET | `/api/users/me` | Get current authenticated user | Completed |
| POST | `/api/users/register` | Register new user | Completed |
| POST | `/api/users/login` | Login user | Completed |
| PUT | TBD | Update user | In progress |
| DELETE | TBD | Delete user | In progress |

### Register User

```http
POST /api/users/register
```

Request body:

```json
{
  "Name": "Dusan Stoiljkovic",
  "Email": "dusan@gmail.com",
  "Password": "dusan123"
}
```

### Login User

```http
POST /api/users/login
```

Request body:

```json
{
  "Email": "dusan@gmail.com",
  "Password": "dusan123"
}
```

Expected response:

```json
{
  "jwt": "your.jwt.token"
}
```

### Get Current User

```http
GET /api/users/me
```

Headers:

```http
Authorization: your.jwt.token
```

### Get User By ID

```http
GET /api/users/{user_id}
```

Example:

```http
GET /api/users/1
```

Headers:

```http
Authorization: your.jwt.token
```

### User Service - In Progress

These requests exist in Postman but do not currently have defined URLs:

| Method | Request Name | Status |
|---|---|---|
| GET | Get All Users | In progress |
| PUT | Update User | In progress |
| DELETE | Delete User | In progress |

---

## Movie Service

## Movies

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | TBD | Search movies | In progress |
| GET | `/api/movies` | Get all movies | Completed |
| GET | `/api/movies/{id}` | Get movie by ID | Completed |
| GET | TBD | Get movies by genre | In progress |
| POST | `/api/movies` | Create movie | Completed |
| PUT | `/api/movies/{id}` | Update movie by ID | Completed |
| DELETE | `/api/{id}` | Delete movie by ID | Completed |

> Note: The Postman collection currently defines `Delete Movie By Id` as `DELETE /api/3`. Verify whether the intended endpoint should be `DELETE /api/movies/{id}`.

### Get All Movies

```http
GET /api/movies
```

### Get Movie By ID

```http
GET /api/movies/{id}
```

Example:

```http
GET /api/movies/1
```

### Create Movie

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
  "genre_ids": [1, 2, 6]
}
```

### Update Movie By ID

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
  "rating": 7.0,
  "genre_ids": [1, 2]
}
```

### Delete Movie By ID

```http
DELETE /api/{id}
```

Example from Postman:

```http
DELETE /api/3
```

### Movies - In Progress

These requests exist in Postman but do not currently have defined URLs:

| Method | Request Name | Status |
|---|---|---|
| GET | Search Movies | In progress |
| GET | Get Movies By Genre | In progress |

---

## Genres

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | `/api/genres` | Get all genres | Completed |
| POST | `/api/genres` | Create genre | Completed |
| DELETE | `/api/movies/genres/{genre_id}` | Delete genre by ID | Completed |

> Note: The Postman collection currently defines `Delete Genre By Id` as `DELETE /api/movies/genres/2`. Verify whether the intended endpoint should be `DELETE /api/genres/{id}`.

### Get All Genres

```http
GET /api/genres
```

### Create Genre

```http
POST /api/genres
```

Request body:

```json
{
  "name": "romance"
}
```

### Delete Genre By ID

```http
DELETE /api/movies/genres/{genre_id}
```

Example from Postman:

```http
DELETE /api/movies/genres/2
```

---

## Movie-Genre Relations

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | `/api/movies/{movie_id}/genres` | Get genres by movie ID | Completed |
| POST | `/api/movies/{movie_id}/genres/{genre_id}` | Create movie-genre relation | Completed |
| DELETE | `/api/movies/{movie_id}/genres/{genre_id}` | Delete movie-genre relation | Completed |

### Get Relation By Movie ID

```http
GET /api/movies/{movie_id}/genres
```

Example:

```http
GET /api/movies/1/genres
```

### Create Relation

```http
POST /api/movies/{movie_id}/genres/{genre_id}
```

Example:

```http
POST /api/movies/1/genres/3
```

### Delete Relation

```http
DELETE /api/movies/{movie_id}/genres/{genre_id}
```

Example:

```http
DELETE /api/movies/1/genres/2
```

---

## Booking Service

## Halls

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | `/api/halls` | Get all halls | Completed |
| GET | TBD | Get projections by hall | In progress |
| GET | `/api/halls/{hall_id}` | Get hall by ID | Completed |
| POST | `/api/halls` | Create hall | Completed |
| PUT | `/api/halls/{hall_id}` | Update hall | Completed |
| DELETE | `/api/halls/{hall_id}` | Delete hall | Completed |

### Get All Halls

```http
GET /api/halls
```

### Get Hall By ID

```http
GET /api/halls/{hall_id}
```

Example:

```http
GET /api/halls/4
```

### Create Hall

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

### Update Hall

```http
PUT /api/halls/{hall_id}
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

### Delete Hall

```http
DELETE /api/halls/{hall_id}
```

Example:

```http
DELETE /api/halls/1
```

### Halls - In Progress

These requests exist in Postman but do not currently have defined URLs:

| Method | Request Name | Status |
|---|---|---|
| GET | Get Projections By Hall | In progress |

---

## Projections

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | `/api/projections` | Get all projections | Completed |
| GET | `/api/projections/{projection_id}` | Get projection by ID | Completed |
| GET | `/api/projections/movie/{movie_id}` | Get projections by movie ID | Completed |
| GET | TBD | Get projections by hall ID | In progress |
| POST | `/api/projections` | Create projection | Completed |
| PUT | `/api/projections/{projection_id}` | Update projection | Completed |
| DELETE | `/api/projections/{projection_id}` | Delete projection | Completed |

### Get All Projections

```http
GET /api/projections
```

### Get Projection By ID

```http
GET /api/projections/{projection_id}
```

Example:

```http
GET /api/projections/1
```

### Get Projections By Movie ID

```http
GET /api/projections/movie/{movie_id}
```

Example:

```http
GET /api/projections/movie/1
```

### Create Projection

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

### Update Projection

```http
PUT /api/projections/{projection_id}
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

### Delete Projection

```http
DELETE /api/projections/{projection_id}
```

Example:

```http
DELETE /api/projections/3
```

### Projections - In Progress

These requests exist in Postman but do not currently have defined URLs:

| Method | Request Name | Status |
|---|---|---|
| GET | Get Projections By HallID | In progress |

---

## Tickets

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | `/api/tickets` | Get all tickets | Completed |
| GET | `/api/tickets/{ticket_id}` | Get ticket by ID | Completed |
| GET | `/api/tickets/users/{user_id}` | Get tickets by user ID | Completed |
| GET | `/api/tickets/projections/{projection_id}` | Get tickets by projection ID | Completed |
| GET | TBD | Get all tickets by order ID | In progress |
| POST | `/api/tickets` | Create ticket | Completed |
| PATCH | TBD | Cancel ticket | In progress |
| DELETE | TBD | Delete ticket | In progress |

### Get All Tickets

```http
GET /api/tickets
```

### Get Ticket By ID

```http
GET /api/tickets/{ticket_id}
```

Example:

```http
GET /api/tickets/1
```

### Get Tickets By User ID

```http
GET /api/tickets/users/{user_id}
```

Example:

```http
GET /api/tickets/users/1
```

### Get Tickets By Projection ID

```http
GET /api/tickets/projections/{projection_id}
```

Example:

```http
GET /api/tickets/projections/1
```

### Create Ticket

```http
POST /api/tickets
```

Request body:

```json
{
  "order_id": 1,
  "projection_id": 1,
  "seat_number": 1
}
```

### Tickets - In Progress

These requests exist in Postman but do not currently have defined URLs:

| Method | Request Name | Status |
|---|---|---|
| GET | Get All Tickets By OrderID | In progress |
| PATCH | Cancel Ticket | In progress |
| DELETE | Delete Ticket | In progress |

---

## Seats

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | TBD | Get reserved seats by projection ID | In progress |
| GET | TBD | Get paid seats by projection ID | In progress |

These requests exist in Postman but do not currently have defined URLs.

---

## Orders

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | TBD | Get all orders | In progress |
| GET | TBD | Get order by ID | In progress |
| GET | TBD | Get orders by user ID | In progress |
| GET | TBD | Get my orders | In progress |
| POST | TBD | Create order | In progress |
| PATCH | TBD | Pay order | In progress |
| PATCH | TBD | Cancel order | In progress |
| DELETE | TBD | Delete order | In progress |

These requests exist in Postman but do not currently have defined URLs.

---

## Recommendation Service

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | TBD | Get recommendations by user ID | In progress |
| GET | TBD | Get my recommendations | In progress |
| GET | TBD | Get similar movies | In progress |

These requests exist in Postman but do not currently have defined URLs.

---

## In Progress Endpoints Summary

The following HTTP requests are not fully defined in the current Postman collection because their URL is missing or empty:

| Service | Method | Request Name | Status |
|---|---|---|---|
| User Service | GET | Get All Users | In progress |
| User Service | PUT | Update User | In progress |
| User Service | DELETE | Delete User | In progress |
| Movie Service / Movies | GET | Search Movies | In progress |
| Movie Service / Movies | GET | Get Movies By Genre | In progress |
| Booking Service / Halls | GET | Get Projections By Hall | In progress |
| Booking Service / Projections | GET | Get Projections By HallID | In progress |
| Booking Service / Tickets | GET | Get All Tickets By OrderID | In progress |
| Booking Service / Tickets | PATCH | Cancel Ticket | In progress |
| Booking Service / Tickets | DELETE | Delete Ticket | In progress |
| Booking Service / Seats | GET | Get Reserved Seats By ProjectionID | In progress |
| Booking Service / Seats | GET | Get Paid Seats By ProjectionID | In progress |
| Booking Service / Orders | GET | Get All Orders | In progress |
| Booking Service / Orders | GET | Get Order By ID | In progress |
| Booking Service / Orders | GET | Get Orders By UserID | In progress |
| Booking Service / Orders | GET | Get My Orders | In progress |
| Booking Service / Orders | POST | Create Order | In progress |
| Booking Service / Orders | PATCH | Pay Order | In progress |
| Booking Service / Orders | PATCH | Cancel Order | In progress |
| Booking Service / Orders | DELETE | Delete Order | In progress |
| Recommendation Service | GET | Get Recommendations By UserID | In progress |
| Recommendation Service | GET | Get My Recommendations | In progress |
| Recommendation Service | GET | Get Similar Movies | In progress |

---

## Error Response Format

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
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource conflict |
| `INTERNAL_ERROR` | 500 | Internal server error |

---

## Authentication

Some routes require JWT authentication.

Login returns a JWT token:

```json
{
  "jwt": "your.jwt.token"
}
```

Use the token in protected routes:

```http
Authorization: Bearer your.jwt.token
```

Do not commit real JWT tokens, passwords, or `.env` files to GitHub.

---

## Development Notes

Run Go formatting in each service:

```bash
gofmt -w .
```

Run a specific service locally:

```bash
go run ./cmd
```

Build Docker images:

```bash
docker compose build
```

Start all services:

```bash
docker compose up
```

---

## Current Features

- User registration
- User login
- JWT authentication
- Get current user
- Get user by ID
- Movie CRUD
- Genre management
- Movie-genre relations
- Hall CRUD
- Projection CRUD
- Ticket creation and ticket queries
- API Gateway routing
- MySQL databases
- Docker Compose setup
- SafeError-style centralized error handling

---

## Future Improvements

- Finish all endpoints marked as `In progress`
- Complete orders flow
- Complete seat availability endpoints
- Complete recommendation-service endpoints
- Add unit tests for service layer
- Add integration tests with test database
- Add database migrations
- Add health check endpoints
- Add CI/CD pipeline
- Add Swagger/OpenAPI documentation
- Improve logging with request IDs
- Add role-based authorization
- Add payment flow for tickets
- Add event-driven communication between services
