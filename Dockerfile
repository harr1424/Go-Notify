FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y awscli

COPY frost /app/frost
WORKDIR /app

CMD ["./frost"]