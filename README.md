# Cinema Microservices System

Cinema Microservices System is a backend project built with Go using a microservices-style architecture.

The system is designed for managing users, movies, genres, halls, projections, and tickets for a cinema platform.

The project contains multiple independent services connected through an API Gateway.

---

## Architecture

The system is divided into several services:

| Service | Responsibility |
|---|---|
| API Gateway | Single entry point for client requests |
| User Service | User registration, login, authentication, and user profile |
| Movie Service | Movies, genres, and movie-genre relations |
| Booking Service | Halls, projections, and tickets |

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
```

Do not commit real `.env` files with secrets to GitHub.

---

## API Gateway

All requests should go through the API Gateway.

Example base URL:

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
```

---

# API Routes

## User Service

### Register User

```http
POST /api/users/register
```

Request body:

```json
{
  "name": "Dusan",
  "email": "dusan@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "name": "Dusan",
  "email": "dusan@example.com",
  "created_at": "2026-05-05T12:00:00Z"
}
```

---

### Login User

```http
POST /api/users/login
```

Request body:

```json
{
  "email": "dusan@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "jwt": "your.jwt.token"
}
```

---

### Get Current User

```http
GET /api/users/me
```

Headers:

```http
Authorization: Bearer your.jwt.token
```

---

### Get User By ID

```http
GET /api/users/{id}
```

Example:

```http
GET /api/users/1
```

---

## Movie Service

## Movies

### Get All Movies

```http
GET /api/movies
```

Optional query parameters:

| Parameter | Description |
|---|---|
| `limit` | Number of movies to return |
| `offset` | Number of movies to skip |
| `sort` | Sorting option |
| `genre` | Filter by genre name |
| `min_year` | Minimum release year |
| `max_year` | Maximum release year |
| `min_rating` | Minimum rating |

Example:

```http
GET /api/movies?limit=10&offset=0&sort=rating_desc&genre=Action&min_rating=7
```

---

### Get Movie By ID

```http
GET /api/movies/{id}
```

Example:

```http
GET /api/movies/1
```

---

### Create Movie

```http
POST /api/movies
```

Request body:

```json
{
  "title": "The Dark Knight",
  "description": "Batman faces the Joker in Gotham City.",
  "year": 2008,
  "image_url": "https://example.com/dark-knight.jpg",
  "duration": 152,
  "rating": 9.0,
  "genres": [
    {
      "id": 1
    },
    {
      "id": 2
    }
  ]
}
```

---

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
  "title": "The Dark Knight",
  "description": "Updated movie description.",
  "year": 2008,
  "image_url": "https://example.com/dark-knight.jpg",
  "duration": 152,
  "rating": 9.1,
  "genres": [
    {
      "id": 1
    }
  ]
}
```

---

### Delete Movie By ID

```http
DELETE /api/movies/{id}
```

Example:

```http
DELETE /api/movies/1
```

---

## Genres

### Get All Genres

```http
GET /api/genres
```

---

### Create Genre

```http
POST /api/genres
```

Request body:

```json
{
  "name": "Action"
}
```

---

### Delete Genre By ID

```http
DELETE /api/genres/{id}
```

Example:

```http
DELETE /api/genres/1
```

---

## Movie-Genre Relations

### Get Relations By Movie ID

```http
GET /api/movies/{id}/genres
```

Example:

```http
GET /api/movies/1/genres
```

---

### Create Relation

```http
POST /api/movies/{movieId}/genres/{genreId}
```

Example:

```http
POST /api/movies/1/genres/2
```

---

### Delete Relation

```http
DELETE /api/movies/{movieId}/genres/{genreId}
```

Example:

```http
DELETE /api/movies/1/genres/2
```

---

## Booking Service

## Halls

### Get All Halls

```http
GET /api/halls
```

---

### Get Hall By ID

```http
GET /api/halls/{id}
```

Example:

```http
GET /api/halls/1
```

---

### Create Hall

```http
POST /api/halls
```

Request body:

```json
{
  "name": "Sala 1",
  "location": "Pozoriste",
  "capacity": 200
}
```

---

### Update Hall

```http
PUT /api/halls/{id}
```

Example:

```http
PUT /api/halls/1
```

Request body:

```json
{
  "name": "Sala 1",
  "location": "Pozoriste",
  "capacity": 250
}
```

---

### Delete Hall

```http
DELETE /api/halls/{id}
```

Example:

```http
DELETE /api/halls/1
```

---

## Projections

### Get All Projections

```http
GET /api/projections
```

---

### Get Projection By ID

```http
GET /api/projections/{id}
```

Example:

```http
GET /api/projections/1
```

---

### Get Projections By Movie

```http
GET /api/projections/movie/{movie_id}
```

Example:

```http
GET /api/projections/movie/1
```

---

### Create Projection

```http
POST /api/projections
```

Request body:

```json
{
  "movie_id": 1,
  "hall_id": 1,
  "start_time": "2026-05-05T18:00:00+02:00",
  "end_time": "2026-05-05T20:00:00+02:00",
  "price": 500
}
```

---

### Update Projection

```http
PUT /api/projections/{id}
```

Example:

```http
PUT /api/projections/1
```

Request body:

```json
{
  "movie_id": 1,
  "hall_id": 1,
  "start_time": "2026-05-05T19:00:00+02:00",
  "end_time": "2026-05-05T21:00:00+02:00",
  "price": 550
}
```

---

### Delete Projection

```http
DELETE /api/projections/{id}
```

Example:

```http
DELETE /api/projections/1
```

---

## Tickets

### Get All Tickets

```http
GET /api/tickets
```

---

### Get Ticket By ID

```http
GET /api/tickets/{id}
```

Example:

```http
GET /api/tickets/1
```

---

### Get Ticket By User

```http
GET /api/tickets/user/{user_id}
```

Example:

```http
GET /api/tickets/user/1
```

---

### Get Ticket By Schedule

```http
GET /api/tickets/schedule/{schedule_id}
```

Example:

```http
GET /api/tickets/schedule/1
```

---

### Create Ticket

```http
POST /api/tickets
```

Request body:

```json
{
  "user_id": 1,
  "projection_id": 1,
  "seat_number": "A10"
}
```

---

### Cancel Ticket

```http
PATCH /api/tickets/{id}/cancel
```

Example:

```http
PATCH /api/tickets/1/cancel
```

---

### Delete Ticket

```http
DELETE /api/tickets/{id}
```

Example:

```http
DELETE /api/tickets/1
```

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
- Movie CRUD
- Genre management
- Movie-genre relations
- Hall CRUD
- Projection CRUD
- Ticket management
- API Gateway routing
- MySQL databases
- Docker Compose setup
- SafeError-style centralized error handling

---

## Future Improvements

- Unit tests for service layer
- Integration tests with test database
- Database migrations
- Health check endpoints
- CI/CD pipeline
- Swagger/OpenAPI documentation update
- Better logging with request IDs
- Role-based authorization
- Payment flow for tickets
- Event-driven communication between services