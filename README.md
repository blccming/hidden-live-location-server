## private positioning server

### TODO
#### storage
- [ ] **use valkey to store location data for all users**
- [ ] **expiration logic in valkey**
  - [ ] delete location data of session after TTL timeout (different namespace from session data)
  - [ ] terminate session after session timeout
#### deployment
- [ ] use docker containers for development / testing and deployment
- [X] use environment variables for log levels and host/port configuration
- [ ] add github actions for CI/CD (build binary and container image)
### optimizations
- [ ] "make secure"
  - [ ] use password authentication?
- [X] add rate limiter per-client (besides the global one)
- [X] add proper logging (with zerolog?)
- [X] file management
- [ ] Use links to share location session => make link directly redirect to app (this is mostly handled client-side)
- [X] Evaluate the usage of Valkey compared to Redis => use Valkey for location data storage

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
