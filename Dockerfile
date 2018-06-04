FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/elahi-arman/urlshortener

ENV SHORTLY_HOME /opt/shortly/
RUN mkdir -p /opt/shortly/config && \
    mkdir /opt/shortly/logs && \
    touch /opt/shortly/logs/access.log && \
    touch /opt/shortly/logs/app.log

COPY ./config/config.yaml /opt/shortly/config/config.yaml

RUN wget -O installDep.sh https://raw.githubusercontent.com/golang/dep/master/install.sh && \
    echo '135b424e4f922c141ae61e31872e34605170fff58d28c1af055dbb84c3a34f1b  installDep.sh' | sha256sum -c - && \
    sh installDep.sh

WORKDIR /go/src/github.com/elahi-arman/urlshortener
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only && \
    go install .

ENTRYPOINT /go/bin/urlshortener
EXPOSE 48290