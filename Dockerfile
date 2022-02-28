FROM ubuntu:18.04

RUN apt upgrade -y

RUN apt-get update -y

WORKDIR /app

COPY . .

EXPOSE 9001

CMD [ "./main" ]