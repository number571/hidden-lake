FROM --platform=linux/amd64 ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.23.0.linux-amd64.tar.gz

WORKDIR /hidden-lake
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hlc ./cmd/hlc

ENV SERVICE_NETWORK=""
ENV SERVICE_PATH="."
ENV SERVICE_THREADS="1"
CMD ./hlc --path "${SERVICE_PATH}" --network "${SERVICE_NETWORK}" --threads ="${SERVICE_THREADS}"
