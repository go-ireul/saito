FROM alpine

RUN apk --no-cache add ca-certificates

ADD saito /usr/local/bin

CMD ["saito"]
