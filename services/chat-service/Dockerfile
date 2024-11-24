FROM alpine:latest

# Set up app directory
RUN mkdir /app
WORKDIR /app

# Copy the binary and configuration files
COPY bin/chat-service /app
COPY config.example.yaml /app
COPY utils/ent/migrate/migrations /app/utils/ent/migrate/migrations
COPY docs /app/docs 
COPY views /app/views         
COPY public /app/public   
COPY server.crt /app
COPY server.key /app

CMD [ "/app/chat-service" ]
