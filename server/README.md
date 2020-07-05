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
docker run -p 8080:8080 dependency-report-server
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

https://server-u2g5mawisa-uw.a.run.app/ping

## Routes

```bash
GET http://localhost:8080/generate-report
```

```bash
GET http://localhost:8080/ping
```

### Test

https://github.com/search?l=JSON&q=org%3Avendasta+%22%40vendasta%2Fcore%22%3A+%22%5E44.15.1%22&type=Code

 "@vendasta/core": "^44.15.1",

============

Avaliando :   "@angular/animations": 