FROM golang:1.16

COPY . /product
WORKDIR /product

# RUN go mod init
# RUN go get github.com/gofiber/fiber
# RUN go get github.com/google/uuid
# RUN go get github.com/gorilla/mux
# RUN go get github.com/smartystreets/goconvey 
# RUN go get go.mongodb.org/mongo-driver

ENTRYPOINT APP_ENV=test go run main.go 