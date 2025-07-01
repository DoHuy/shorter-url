# Shorter REST API

A simple URL shortener REST API built with Go and Gin.

## Features

- Create short URLs for any original URL
- Redirect to the original URL using the short code
- Retrieve short URL details by code
- Swagger/OpenAPI documentation

## Requirements

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Devbox](https://www.jetpack.io/devbox/)
- [Make](https://www.gnu.org/software/make/)
- [Go](https://go.dev/) (if running locally, outside Docker)

## Getting Started

### Clone the repository

```sh
git clone https://github.com/yourusername/shorter-rest-api.git
cd shorter-rest-api
```

### Development Environment

You can use [Devbox](https://www.jetpack.io/devbox/) to set up your development environment:

```sh
devbox shell
```

### Build and Run

Build the project:

```sh
make build
```

Run the project:

```sh
make run
```

### Run Tests

```sh
make test
```

### Using Docker

Build and run the app with Docker:

```sh
docker build -t shorter-rest-api .
docker run -p 8080:8080 shorter-rest-api
```

## API Documentation

After running the app, access Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License
