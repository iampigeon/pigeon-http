FROM golang:alpine

COPY /bin/pigeon-http /
WORKDIR /
EXPOSE 9020

CMD ["./pigeon-http"]
