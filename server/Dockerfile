FROM golang:1.20-alpine as go-builder

WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /tcp-client-server /app/server/server.go

FROM scratch as go-runtime-container

COPY --from=go-builder /tcp-client-server ./tcp-client-server

EXPOSE 8080

CMD [ "./tcp-client-server"]