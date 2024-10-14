# Install and run
1. Copy app.env.example as app.env
2. Copy docker/.env.example as docker/.env
3. Build docker containers and run server:
`make up`

# Swagger
Once server is up, swagger docs will be at `http://localhost:8080/swagger/index.html`

# Stop server
To stop server:
`make down`