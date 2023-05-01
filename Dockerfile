ARG GO_VERSION=1.18

FROM golang:${GO_VERSION}-alpine AS builder

# Compile file
# We aren't using a proxy, to use dependencies directly from go module
RUN go env -w GOPROXY=direct
# We need git to install the dependencies
RUN apk add --no-cache git
# Adding and updating security certificates
RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /src

# We need main.go to install the dependencies, otherwhise go mod vendor will return an error
COPY ./main.go ./go.mod ./go.sum ./

RUN go mod download

# Copying all go directories
COPY util util
COPY controller controller
COPY models models

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /gounittesting

FROM scratch AS runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /gounittesting /go-unit-testing

EXPOSE 5000

ENTRYPOINT [ "/go-unit-testing" ]
