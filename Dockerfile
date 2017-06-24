FROM golang:1.8

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/jeremyroberts0/kirk
WORKDIR /go/src/github.com/jeremyroberts0/kirk

RUN go get ./
RUN go build

CMD if [ ${APP_ENV} != development ]; \
	then \
	kirk; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi