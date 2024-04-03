FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./bin/featSvc /app/featSvc

RUN chmod +x /app/featSvc

CMD ["./featSvc"]

