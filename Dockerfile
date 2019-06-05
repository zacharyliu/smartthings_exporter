FROM golang:1.11
WORKDIR /go/src/github.com/zacharyliu/smartthings_exporter
COPY . .
RUN go get .
ENTRYPOINT smartthings_exporter --smartthings.oauth-client=$SMARTTHINGS_OAUTH_CLIENT --smartthings.oauth-token.file=$SMARTTHINGS_OAUTH_TOKEN_FILE
