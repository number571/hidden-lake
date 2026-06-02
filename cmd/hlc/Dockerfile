FROM --platform=linux/amd64 ubuntu:20.04

RUN apt-get update && apt-get install -y \
    wget \
    gcc \
    build-essential \
    libssl-dev \
    zlib1g-dev \
    libncurses5-dev \
    libgdbm-dev \
    libnss3-dev \
    libsqlite3-dev \
    libreadline-dev \
    libffi-dev \
    curl \
    libbz2-dev \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://www.python.org/ftp/python/3.14.5/Python-3.14.5.tgz && \
    tar -C /tmp -xzf Python-3.14.5.tgz

WORKDIR /tmp/Python-3.14.5
RUN ./configure --enable-optimizations && make -j$(nproc) && make install

RUN ln -sf /usr/local/bin/python3.14 /usr/local/bin/python && \
    ln -sf /usr/local/bin/pip3.14 /usr/local/bin/pip && \
    rm -rf /tmp/Python-3.14.5*

RUN wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz && \ 
    tar -C /opt -xzf go1.23.0.linux-amd64.tar.gz

WORKDIR /hidden-lake
ENV PATH="${PATH}:/opt/go/bin"

COPY ./ ./
RUN go build -o hlc ./cmd/hlc

ENV SERVICE_NETWORK=""
ENV SERVICE_PATH="."
CMD ./hlc --path "${SERVICE_PATH}" --network "${SERVICE_NETWORK}"
