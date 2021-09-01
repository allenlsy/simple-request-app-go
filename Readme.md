Simple Request App
---

This simple application sends request to a list of endpoints, and also listens to incoming request.

The configuration of the application is based on environment variables:

* `POD_NAME`: name of the kubernetes pod. It can also be used as host name.
* `LISTEN_PORT`: port number that server listens to incoming requests
* `BACKEND_PORT`: port number that server sends outgoing request to
* `SERVICE_NAME`: name of the service that application belongs to. It will also be the response message to all incoming requests
* `BACKEND_ENDPOINTS`: a comma-separated list of endpoints. E.g. `"google.com,youtube.com"`
* `REQUEST_RATE`: number of outgoing requests per second. It can be set to integer or float
* `LIFETIME`: time that pod will be alive. If not provided, the pod will live forever. Time format are like `"15m"`, `"60s"`. A less than 30 seconds random duration will be added to lifetime.

### Run the application

Start server in one terminal:

```
BACKEND_ENDPOINTS="www.google.com" LISTEN_PORT=8080 REQUEST_RATE=0.5 SERVICE_NAME=my_service POD_NAME=abc123 go run cmd/client-app/main.go
```

Query server in another terminal. Here uses [hey](https://github.com/rakyll/hey) to send the requests:

```
# -q: QPS per worker
# -c: number of workers
hey -q 2 -c 1 http://localhost:8080/
```