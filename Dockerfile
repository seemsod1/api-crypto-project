FROM golang

WORKDIR /app

COPY . /app

RUN go get .
RUN go build -v -o bin .


CMD ["/app/bin"]