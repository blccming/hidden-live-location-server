## hidden live location server

### TODO
#### storage
- [X] **use valkey to store location data for all users**
- [X] **expiration logic in valkey**
  - [X] delete location data of session after TTL timeout (different namespace from session data)
  - [X] terminate session after session timeout
- [ ] use valkey api in endpoints.go
#### deployment
- [ ] use docker containers for development / testing and deployment
- [ ] add valkey configuration parameters to env vars
- [ ] add github actions for CI/CD (build binary and container image)
### optimizations
- [ ] "make secure"
  - [ ] use password authentication?
- [ ] Use links to share location session => make link directly redirect to app (this is mostly handled client-side)

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
