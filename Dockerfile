FROM golang:1.8

ARG app_env
ENV APP_ENV $app_env

COPY ./ /go/src/github.com/ContainX/kirk
WORKDIR /go/src/github.com/ContainX/kirk

RUN go get ./
RUN go build

EXPOSE 8080:8080

CMD if [ ${APP_ENV} != development ]; \
	then \
	kirk; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi