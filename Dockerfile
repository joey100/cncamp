FROM golang:1.17 AS build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN mkdir src/httpserver
WORKDIR src/httpserver
COPY ./httpserver ./
# If don't use glog module, then we don't need below
RUN go mod init && go mod tidy
RUN go build -o /bin/httpserver

FROM scratch
COPY --from=build /bin/httpserver /bin/httpserver
ENTRYPOINT ["/bin/httpserver"]
