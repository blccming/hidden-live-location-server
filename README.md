## private positioning server

### TODO
#### storage
- [ ] use redis to store location data for all users [TOKEN, LONGITUDE, LATIDUDE, TTL, TIMEOUT, LAST_UPDATE (, PW-HASH?)]
- [ ] coroutine logic
  - [ ] delete data of session after TTL timeout
  - [ ] termiante session after session timeout
#### deployment
- [ ] use docker containers for development / testing and deployment
- [X] use environment variables for log levels and host/port configuration
- [ ] add github actions for CI/CD (build binary and container image)
### optimizations
- [ ] "make secure"
  - [ ] use password authentication?
- [ ] **add rate limiter per-client (besides the global one)**
- [ ] **add proper logging (with zerolog?)**
- [X] file management
- [ ] Use links to share location session => make link directly redirect to app (this is mostly handled client-side)
- [ ] Evaluate the usage of Valkey compared to Redis

### Development
Install docker with docker compose if not already installed.
```sh
curl -fsSL https://get.docker.com | sh
```

Configure your Cloudflare Zero Trust Token in `docker/.env` and start the container.
```sh
cd docker
mv .env.example .env
nano .env # replace placeholder with your token
docker compose up -d
```

Now the project can be executed while being available via your domain.
```sh
go run .
```

When making changes to the API, document them and update swagger docs (make sure `swag` is in PATH) with
```sh
swag init
```

API docs are available at `your-host.tld/docs/index.html` (when hosting on your machine at `localhost:8080/docs/index.html`).
