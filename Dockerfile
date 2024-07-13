FROM ubuntu:jammy-20240530

ENV DEBIAN_FRONTEND=noninteractive
RUN  apt-get update &&  apt-get install -y awscli

COPY frost /app/frost
COPY apnkey.p8 /app/apnkey.p8

WORKDIR /app

CMD ["./frost"]