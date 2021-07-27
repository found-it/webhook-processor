FROM golang:1.14-alpine as builder

ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOBIN
ENV PATH=$PATH:$GOROOT/bin

# Install the application
WORKDIR $GOPATH/src/webhook
COPY . ./

RUN apk update && apk add --no-cache git

RUN go build -o $GOBIN/webhook

EXPOSE 9000

USER foundit

FROM alpine:latest
COPY --from=builder /go/bin/webhook /bin/webhook
RUN addgroup --gid 2323 "foundit" && \
    adduser --disabled-password \
            --home "/home/foundit" \
            --ingroup "foundit" \
            --no-create-home \
            --uid 2324 \
            "foundit"

USER foundit
ENTRYPOINT ["webhook"]

