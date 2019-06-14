FROM golang:1.9.2
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/go-sql-driver/mysql"
WORKDIR $GOPATH/src/main
ADD . $GOPATH/src/main
RUN go build .
EXPOSE 8080
ENTRYPOINT ["./main"]