# Dependency Report - Server

Requires Go 1.27+.

```bash
GITHUB_AUTH_TOKEN=<token> go run .
GITHUB_AUTH_TOKEN=$(gh auth token) go run .
```

## Tests

```bash
go test ./...
```

## Docker

### Build the image

```bash
docker build . -t dependency-report-server
```

### Run the image

`GITHUB_AUTH_TOKEN` is not baked into the image, it must be supplied at runtime:

```bash
docker run -p 8080:8080 -e GITHUB_AUTH_TOKEN=<token> dependency-report-server
```

## Cloud Run

### Submit new image

```bash
gcloud builds submit --tag gcr.io/dependency-report/server
```

### Deploy the new image

```bash
gcloud run deploy --image gcr.io/dependency-report/server --platform managed
```

### Server production

---

## Routes

```bash
GET http://localhost:8080/generate-report
```

```bash
GET http://localhost:8080/ping
```

## Libraries

- https://github.com/gofiber/fiber
- https://github.com/google/go-github
- https://github.com/awslabs/git-secrets

### Upgrade Go Fiber

```bash
go get -u github.com/gofiber/fiber/v3
go get -u github.com/gofiber/fiber/v3/middleware/cors
go get -u github.com/gofiber/fiber/v3/middleware/logger
```