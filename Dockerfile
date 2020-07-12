FROM golang:alpine
ENV GO111MODULE=on
RUN mkdir /app 
ADD . /app/
WORKDIR /app 

RUN go mod download

RUN go build -o main .
CMD ["./main"]