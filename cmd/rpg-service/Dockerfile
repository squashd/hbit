FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./bin/rpgSvc /app/rpgSvc

RUN chmod +x /app/rpgSvc

CMD ["./rpgSvc"]

