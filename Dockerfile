FROM golang:alpine AS server-builder

RUN apk --no-cache add git

RUN go get github.com/danielwoodsdeveloper/offline_coding_challenges/server/runtimes
RUN go get github.com/danielwoodsdeveloper/offline_coding_challenges/server/tests
RUN go get github.com/docker/docker/api/types
RUN go get github.com/docker/docker/api/types/container
RUN go get github.com/docker/docker/api/types/mount
RUN go get github.com/docker/docker/client
RUN go get github.com/gorilla/handlers
RUN go get github.com/gorilla/mux
RUN go get github.com/icza/gox/stringsx
RUN go get github.com/sony/sonyflake

COPY server .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server

FROM node:alpine AS client-builder

WORKDIR /app

COPY client/coding-challenges .

RUN npm install
RUN npm install react-scripts

FROM alpine

COPY --from=server-builder /server /server

RUN apk --no-cache add docker

RUN mkdir -p /static
COPY --from=client-builder /app/build /static

RUN chmod +x /server

CMD [ /server ]