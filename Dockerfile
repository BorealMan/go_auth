FROM golang:alpine3.18

WORKDIR /app

RUN apk add --no-cache git

RUN git clone https://github.com/BorealMan/go_fiber_mysqlx_user_auth.git

RUN cd /app/go_fiber_mysqlx_user_auth; go build main.go

ENTRYPOINT [ "/app/go_fiber_mysqlx_user_auth/main" ]

EXPOSE 8000

