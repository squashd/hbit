FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./bin/apiGateway /app/apiGateway

RUN chmod +x /app/apiGateway

CMD ["./apiGateway"]


