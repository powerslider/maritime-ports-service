FROM golang:1.20-alpine as builder

ENV PROJECT_NAME maritime-ports-service
ENV BASE_DIR /go/src/github.com/powerslider/${PROJECT_NAME}
WORKDIR ${BASE_DIR}

RUN apk --no-cache add git ca-certificates

COPY go.mod go.sum ${BASE_DIR}/

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 go build -v -o /dist/${PROJECT_NAME} ./cmd/${PROJECT_NAME}/main.go

FROM alpine

ENV PROJECT_NAME maritime-ports-service
ENV BASE_DIR /go/src/github.com/powerslider/${PROJECT_NAME}

RUN apk --no-cache add ca-certificates

COPY --from=builder /dist .
COPY --from=builder ${BASE_DIR}/fixtures fixtures

CMD ["/maritime-ports-service"]
