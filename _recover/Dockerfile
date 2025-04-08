FROM --platform=linux/amd64 ubuntu:latest
LABEL authors="m41den"

RUN apt update && apt install -y mariadb-client redis-tools wget curl
RUN wget https://dl.min.io/client/mc/release/linux-amd64/mc -O /usr/local/bin/mc && chmod +x /usr/local/bin/mc

WORKDIR /recover
COPY recover.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
