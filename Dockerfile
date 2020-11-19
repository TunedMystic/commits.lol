# ----------------------------------------------------
# In stage one, we're installing dependencies
# and building the golang application.
# ----------------------------------------------------

FROM golang:1.15.5-alpine3.12 as builder

ENV GO111MODULE=on \
    GOOS=linux \
    CGO_ENABLED=1 \
    BUILD_PATH=/build

WORKDIR $BUILD_PATH

RUN apk add --no-cache musl-dev gcc

ADD go.mod go.sum $BUILD_PATH/
RUN go mod download

COPY . $BUILD_PATH
RUN go build -ldflags="-s -w" -o commits.lol



# ----------------------------------------------------
# In stage two, we're copying the app binary and
# other dependencies into a minimal image.
# ----------------------------------------------------

FROM alpine:3.12

ENV APP_PATH=/usr/src

WORKDIR $APP_PATH
RUN apk add --no-cache bash

COPY --from=builder /build/commits.lol $APP_PATH
COPY --from=builder /build/static $APP_PATH/static
COPY --from=builder /build/templates $APP_PATH/templates

EXPOSE 8000
CMD ["./commits.lol", "server"]
