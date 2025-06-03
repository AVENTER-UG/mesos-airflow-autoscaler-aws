FROM golang:alpine AS builder

WORKDIR /build

COPY . /build/

RUN apk update && apk upgrade && apk add git && \
    go get -d

ARG TAG
ARG BUILDDATE
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.BuildVersion=$BUILDDATE -X main.GitVersion=$TAG -extldflags \"-static\"" -o main .

FROM golang:alpine
LABEL maintainer="Andreas Peters <support@aventer.biz>"
LABEL org.opencontainers.image.title="mesos-airflow-autoscaler"
LABEL org.opencontainers.image.description="Airflow Autoscaler for AWS and Apache Mesos/ClusterD"
LABEL org.opencontainers.image.vendor="AVENTER UG (haftungsbeschränkt)"
LABEL org.opencontainers.image.source="https://github.com/AVENTER-UG/"

RUN apk add --no-cache ca-certificates
RUN adduser -S -D -H -h /app appuser

USER appuser

COPY --from=builder /build/main /app/

WORKDIR "/app"

CMD ["./main"]
