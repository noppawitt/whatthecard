FROM golang:1.14.6-alpine3.12 AS go-builder
WORKDIR /go/src/whatthecard
COPY go.mod go.sum ./
COPY main.go ./
COPY pkg pkg
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o whatthecard .

FROM node:12.18.2-alpine3.12 as js-builder
WORKDIR /app
COPY web/package*.json ./
RUN yarn install
COPY web .
RUN yarn build

FROM alpine:3.12.0
WORKDIR /app
COPY --from=js-builder /app/dist ./web/dist
COPY --from=go-builder /go/src/whatthecard/whatthecard ./
CMD ["./whatthecard"]
