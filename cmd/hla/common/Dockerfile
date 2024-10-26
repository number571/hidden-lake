FROM --platform=linux/amd64 ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.23.0.linux-amd64.tar.gz

WORKDIR /hidden-lake
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o common ./cmd/hla/common

ENV SERVICE_NETWORK=""
ENV SERVICE_PATH="/mounted"
CMD ./common -path="${SERVICE_PATH}" -network="${SERVICE_NETWORK}"
