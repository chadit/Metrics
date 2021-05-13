# build binary
FROM golang as builder
COPY . /go/src/github.com/chadit/Metrics
WORKDIR /go/src/github.com/chadit/Metrics
ARG VERSION=0.0.0
ENV TAG_VERSION=$VERSION

RUN CGO_ENABLED=0 go build -a -ldflags "-X main.Build=$VERSION-`date -u +.%Y%m%d.%H%M%S`" -o ./bin/metrics ./cmd/metrics

FROM node:10-alpine

COPY ./bin /go/src/github.com/chadit/Metrics
WORKDIR /go/src/github.com/chadit/Metrics

COPY --from=builder /go/src/github.com/chadit/Metrics/bin/metrics /go/src/github.com/chadit/Metrics/metrics
RUN ls /go/src/github.com/chadit/Metrics -ll

CMD [ "/go/src/github.com/chadit/Metrics/metrics"]