Cinema Microservices System

Cinema Microservices System is a backend project built in Go during my backend internship.  
The goal of the project is to simulate a distributed cinema system using multiple services, REST APIs, Docker and separate databases.

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

Main Features

- User registration and login
- JWT authentication
- Movie management
- Cinema schedule management
- Ticket / booking creation
- API Gateway for routing requests
- Docker Compose setup
- OpenAPI documentation

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

Architecture

The system is split into several backend services:

    | Service | Responsibility |
    |---|
    | API Gateway | Entry point for client requests |
    | User Service | Registration, login, JWT, user data |
    | Movie Service | Movies, genres, schedules |
    | Booking Service | Tickets and reservations |

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

Tech Stack

- Go
- net/http / chi router
- MySQL
- GORM
- Docker
- Docker Compose
- JWT
- OpenAPI

--------------------------------------------------------------------------------------------------------------------------------------------------------------------------

How to Run

```bash
docker compose up --build
