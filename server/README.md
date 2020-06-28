# Dependency Report - Server

```bash
GITHUB_AUTH_TOKEN=<token> go run *.go
```

## Docker

### Build the image

```bash
docker build . -t dependency-report-server
```

### Run the image

```bash
docker run -p 3000:3000 dependency-report-server
```

## Routes

```bash
GET http://localhost:3000/generate-report
```

```bash
GET http://localhost:3000/ping
```

### Test

https://github.com/search?l=JSON&q=org%3Avendasta+%22%40vendasta%2Fcore%22%3A+%22%5E44.15.1%22&type=Code

 "@vendasta/core": "^44.15.1",

============

Avaliando :   "@angular/animations": 