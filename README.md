## private positioning server

### TODO
#### storage
- [ ] use redis to store location data for all users [TOKEN, COORDINATES, LAST_UPDATE]
- [ ] coroutine logic
  - [ ] delete data of session after TTL timeout
  - [ ] termiante session after session timeout
#### deployment
- [ ] use docker containers for development / testing and deployment
### optimizations
- [ ] "make secure"
  - [ ] use password authentication?
- [ ] change rate limiter from global to per-client

### Notes
- Use links to share location session => make link directly redirect to app (through some href? no idea how that works from backend / web server)
- Would it be worth it to use Valkey instead of Redis?
- Be aware of directory traversal attacks (could already be protected by gin, not sure about that though)

### Development
Install docker with docker compose if not already installed.
```sh
curl -fsSL https://get.docker.com | sh
```

Configure your Cloudflare Zero Trust Token in `docker/.env` and start the container.
```sh
cd docker
mv .env.example .env
nano .env
docker compose up -d
```

Now the project can be executed while being available via your domain.
```sh
go run .
```

When making changes to the API, document them and update swagger docs with
```sh
swag init
```

API docs are available at `/docs/index.html`.
