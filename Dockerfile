FROM golang:1.12-alpine as builder
RUN apk add git
COPY . /go/src/sbdb-teach-course
ENV GO111MODULE on
WORKDIR /go/src/sbdb-teach-course
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
COPY --from=builder /go/src/sbdb-teach-course/sbdb-teach-course /
WORKDIR /
CMD ./sbdb-teach-course
ENV PORT 8000
EXPOSE 8000