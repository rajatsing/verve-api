Pre-requisites
----------------------
1. Docker
2. Docker Compose
3. golang 1.20(if you want to run the app locally)


Build and Run the App
----------------------
1. Clone the repository
2. in the project directory run `docker-compose up --build`


Testing The App
----------------------
1. Open a new terminal
2. Without `endpoint` parameter
```bash
curl "http://localhost:8080/api/verve/accept?id=123"
```

1. With `endpoint` parameter
```bash
curl "http://localhost:8080/api/verve/accept?id=123&endpoint=http://httpbin.org/post"
```

Expected Response for both cases
```json
ok
```