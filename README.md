## hidden live location server

### TODO
#### storage
#### deployment
- [ ] add valkey maxmem configuration to env vars
- [ ] add docker images for arm64

### optimizations
- [ ] "make secure"
  - [ ] use password authentication?
- [ ] Use links to share location session => make link directly redirect to app (this is mostly handled client-side)
- [ ] test token generation / randomness -> only one session at a time


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
