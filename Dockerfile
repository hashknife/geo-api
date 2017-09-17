FROM golang:1.9 as builder
COPY ./ /go/src/github.com/hashknife/geo-api/
RUN go get -u -v github.com/go-kit/kit/... && \
    go get -u -v github.com/gorilla/mux && \
    go get -u -v github.com/gorilla/handlers && \
    go get -u -v github.com/garyburd/redigo/... && \
    go get -u -v github.com/stretchr/testify/suite
RUN cd /go/src/github.com/briandowns/hashknife/geo-api/ && \
    go build -o geo-api -v .
FROM centurylink/ca-certs
WORKDIR /
EXPOSE 9988
EXPOSE 9989
COPY --from=builder /gocode/src/github.com/hashknife/geo-api/bin/geo-api .
COPY config.json /
ENTRYPOINT ["/geo-api", "-c", "config.json"]
