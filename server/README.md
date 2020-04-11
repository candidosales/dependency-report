# Server

```bash
GITHUB_AUTH_TOKEN=<token> APP_ENV=development go run *.go
```

## Build the image

```bash
docker build . -t dependency-report-server
```

## Run the image

```bash
docker run -p 3000:3000 dependency-report-server
```