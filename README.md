## private positioning server

### TODO
#### API endpoints
- [ ] make all planned endpoints available
    - [ ] /session
        - [ ] POST /create
            - INPUT: TTL, DURATION (automatic termination at certain timestamp)
            - OUTPUT: JWT (could use golang-jwt/jwt)
        - [ ] POST /terminate
            - INPUT: TOKEN
            - OUTPUT: OK/NOT OK -> could be done via HTTP status
        - [ ] GET /session/<TOKEN>
            - OUTPUT: STATUS (OK/NOT EXISTING/TERMINATED) -> store termination status for x amount, CORDS, LAST_UPDATED
        - [ ] POST /session/<TOKEN>
            - INPUT: CORDS
            - OUTPUT: OK/NOT OK -> could be done via HTTP status
    - [X] GET /health
        - OUTPUT: "OK", RUNTIME
- [ ] document with swagger
#### storage
- [ ] use redis to store location data for all users [TOKEN, COORDINATES, LAST_UPDATE]
#### deployment
- [ ] use docker containers for development / testing and deployment

### Notes
- Could use cloudflare tunnel to make testing deployment available from the internet
- Anti-bruteforce protection: Rate limiting? IP banning (storing IPs would be against design philosophy)?
- Use secrets?
- Use a-z, 0-9 for tokens: e.g. 6 symbols would be 36^4 would be about 1.6 million combinations or 36^6 would be about 217 million combinations => more than enough for our use case
- Use links to share location session => make link directly redirect to app (through some href? no idea how that works from backend / web server)
- Would it be worth it to use Valkey instead of Redis?
- Be aware of directory traversal attacks (could already be protected by gin, not sure about that though)
