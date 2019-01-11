FROM golang:1.11 AS GOLANG_TOOL_CHAIN
WORKDIR /service
COPY . .
RUN make all

FROM alpine:latest
WORKDIR /service/bin
COPY --from=GOLANG_TOOL_CHAIN /service/bin/main /service/bin/main
RUN apk update && apk add ca-certificates
RUN apk add unixodbc-dev tmux vim

EXPOSE 4208

CMD ["/service/bin/main"]
