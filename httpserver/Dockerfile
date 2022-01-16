FROM golang:1.13 AS build
RUN mkdir src/httpserver
WORKDIR src/httpserver
COPY main.go .
RUN go mod init && go mod tidy
RUN go build -o /bin/httpserver

FROM ubuntu
COPY --from=build /bin/httpserver /bin/httpserver
ENTRYPOINT ["/bin/httpserver"]
