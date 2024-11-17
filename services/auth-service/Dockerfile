FROM alpine:latest

# Set up app directory
RUN mkdir /app
WORKDIR /app

# Copy the binary and configuration files
COPY bin/auth-service /app
COPY config.example.yaml /app
COPY utils/ent/migrate/migrations /app/utils/ent/migrate/migrations
COPY docs /app/docs 

CMD [ "/app/auth-service" ]
