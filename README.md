# Movie Theater Backend API

This repository contains the backend API for a movie theater application. 

The application manages movie schedules and movie details, and also offers the ability to make online reservations for tickets. 

The backend is developed using the Go programming language (version go1.22.4) with the Gin framework and MongoDB as the database.

## Features

- Manage movie schedules and details
- Online ticket reservations
- User authentication and authorization
- Secure API endpoints with JWT and PASETO access tokens
- Configuration management using Viper
- Robust unit tests with Testify
- Mocking database interactions for testing

## Technology Stack

* using [go1.22](https://tip.golang.org/doc/go1.22)
* using [gin-gonic](https://github.com/gin-gonic/gin#gin-web-framework) v1.10.0 web framework
* using [viper](https://github.com/spf13/viper) as a configuration solution
* using [mongo-db](https://www.mongodb.com/) as NoSQL DB
* using [jwt-go](github.com/dgrijalva/jwt-go) to provide an implementation of JWT
* using [x/crypto](golang.org/x/crypto), Go Cryptography package 
* using [testify](https://github.com/stretchr/testify), for write robust unit tests 
* using [Gomock](https://github.com/golang/mock) for mocking database


## Getting Started

### Prerequisites
-	Go (version go1.22.4 or higher)
-	MongoDB
-	Git

### Installation
1.	Clone the repository:
```sh
git clone https://github.com/yourusername/movietheater-backend.git
cd movietheater-backend
```
2. Install dependencies:
```sh
go mod download
```
3. Install MongoDB
   
    Instructions for installation are https://www.mongodb.com/docs/manual/installation/

4. Set up environment variables