
FROM golang:latest

WORKDIR /app

ADD . .

#RUN go install 
#RUN go build 
ENTRYPOINT ["tail", "-f", "/dev/null"]
