FROM --platform=linux/amd64 ubuntu:20.04

RUN apt-get update && apt-get install -y wget gcc
RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.23.0.linux-amd64.tar.gz

WORKDIR /hidden-lake
ENV PATH="${PATH}:/opt/go/bin"
COPY ./ ./
RUN go build -o hls ./cmd/hls

ENV SERVICE_PATH="/mounted"
ENV SERVICE_PARALLEL="1"
CMD ./hls -path="${SERVICE_PATH}" -parallel="${SERVICE_PARALLEL}"
