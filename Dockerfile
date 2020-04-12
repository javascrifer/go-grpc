FROM alpine:latest

EXPOSE 50051

WORKDIR /grpc-server

ADD /bin/grpc-server/server ./server

CMD ["./server"]
