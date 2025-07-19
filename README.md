# Backend for a Board project

Realtime board application to draw and watch

Swagger URL: [http://localhost:8080/swagger-ui/index.html](http://localhost:8080/swagger-ui/index.html)

### how to run:

```shell
docker compose up -d
```

### how to build

- container
```shell
docker buildx build . -t ghcr.io/ttodoshi/board-project:latest
```

- locally
```shell
make build
```
