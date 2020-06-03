FROM golang:alpine AS server-builder

RUN apk --no-cache add git

RUN go get github.com/danielwoodsdeveloper/offline_coding_challenges/server/config
RUN go get github.com/danielwoodsdeveloper/offline_coding_challenges/server/tests
RUN go get github.com/docker/docker/api/types
RUN go get github.com/docker/docker/api/types/container
RUN go get github.com/docker/docker/api/types/mount
RUN go get github.com/docker/docker/client
RUN go get github.com/gorilla/mux
RUN go get github.com/sony/sonyflake

COPY server .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server

FROM alpine

COPY --from=server-builder /server /server

RUN apk --no-cache add docker

RUN mkdir -p /static

RUN chmod +x /server

CMD [ /server ]