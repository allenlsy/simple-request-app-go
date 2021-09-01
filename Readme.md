Simple Request App
---

This simple application sends request to a list of endpoints, and also listens to incoming request.

The configuration of the application is based on environment variables:

* `LISTEN_PORT`: port number that server listens to incoming requests
* `BACKEND_PORT`: port number that server sends outgoing request to
* `SERVER_NAME`: name of the server that application runs on. It will also be the response message to all incoming requests
* `BACKEND_ENDPOINTS`: a comma-separated list of endpoints. E.g. `"google.com,youtube.com"`
* `REQUEST_RATE`: number of outgoing requests per second. It can be set to integer or float


### Run the application

Start server in one terminal:

```
BACKEND_ENDPOINTS="www.google.com" LISTEN_PORT=8080 REQUEST_RATE=0.5 SERVER_NAME=my_server go run cmd/client-app/main.go
```

Query server in another terminal. Here uses [hey](https://github.com/rakyll/hey) to send the requests:

```
# -q: QPS per worker
# -c: number of workers
hey -q 2 -c 1 http://localhost:8080/
```