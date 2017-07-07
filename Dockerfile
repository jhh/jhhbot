# docker build -t j3ff/jhhbot .

# docker run --name jhhbot               \
#   --detach --restart unless-stopped    \
#   --env TWITTER_CONSUMER_KEY=""        \
#   --env TWITTER_CONSUMER_SECRET=""     \
#   --env TWITTER_ACCESS_TOKEN=""        \
#   --env TWITTER_ACCESS_TOKEN_SECRET="" \
#   j3ff/jhhbot                                                                                                               

FROM golang:1.8
CMD mkdir -p /go/src/app
COPY . /go/src/app
WORKDIR /go/src/app
RUN go install
CMD ["/go/bin/app"]