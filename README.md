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

1.	**Clone the repository:**

```sh
git clone https://github.com/yourusername/movietheater-backend.git
cd movietheater-backend
```

2. **Install dependencies:**

```sh
go mod download
```

3. **Install MongoDB**
   
Instructions for installation are https://www.mongodb.com/docs/manual/installation/

4. **Set up environment variables**

Create a `.env` file in the root directory of the project and add the necessary configuration variables.

```sh
TOKEN_SYMMETRIC_KEY=your_jwt_secret
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=24h
HTTP_SERVER_ADDRESS=0.0.0.0:8080
MONGO_URL=mongodb://localhost:27017
USERNAME=db_username
PASSWORD=db_password
DATABASE=db_name
JWT_SECRET=your_jwt_secret
```
5. **Start local MongoDB with replica sets:**

Create folder `data` and 3 subfolders: `rs0-0`, `rs0-1`, `rs0-2`.
Inside the `mongosh` shell, run the following command:
```sh
rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "localhost:27017" }
  ]
});
```
After initiating the replica set, add the other members:
```sh
rs.add("localhost:27018");
rs.add("localhost:27019");
rs.status();
```
Open separate command prompt windows (or other shells) and run the following commands to start each replica set:
```sh
mongod --port 27017 --dbpath folder_path\data\rs0-0 --replSet rs0
mongod --port 27018 --dbpath folder_path\data\rs0-1 --replSet rs0
mongod --port 27019 --dbpath folder_path\data\rs0-2 --replSet rs0
```
In another command prompt, connect to the first MongoDB instance and initiate the replica set:
```sh
mongosh --port 27017
```

## Running the Application

1. **Start MongoDB:**

Ensure MongoDB is running on your local machine with the replica set configured as described above.

2. **Run the server:**

```sh
go run main.go
```

3. **Access the API:**

The API will be available at http://localhost:8080.

## Running Tests

To run tests, use the following command:

```sh
go test ./...
```

This will run all the tests in the project, including unit tests with mocked database interactions.

## API Documentation

Detailed documentation and descriptions of the API endpoints are available at:

http://localhost:8080/swagger/index.html

## License

[MIT](https://choosealicense.com/licenses/mit/)
