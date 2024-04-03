FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./bin/taskSvc /app/taskSvc

RUN chmod +x /app/taskSvc

CMD ["./taskSvc"]

