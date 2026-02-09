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
    - [ ] GET /health
        - OUTPUT: "OK", RUNTIME
- [ ] document with swagger
#### storage
- [ ] use redis to store location data for all users [TOKEN, COORDINATES, LAST_UPDATE]
#### deployment
- [ ] use docker containers for development / testing and deployment
