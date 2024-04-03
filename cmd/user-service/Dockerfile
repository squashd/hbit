FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./bin/userSvc /app/userSvc

RUN chmod +x /app/userSvc

CMD ["./userSvc"]

