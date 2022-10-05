FROM golang:1.19

WORKDIR /product-scrapper

COPY Makefile /product-scrapper
RUN make setup

COPY go.mod /product-scrapper
COPY go.sum /product-scrapper
RUN make download

COPY . /product-scrapper
RUN make build

FROM debian:buster-slim
EXPOSE 8888

RUN set -x && \
    apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates

WORKDIR /home/application/current
COPY --from=0 /product-scrapper/dist/api .
RUN chmod +x api
CMD ["./api"]