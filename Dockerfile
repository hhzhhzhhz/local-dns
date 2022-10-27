FROM golang:1.17
WORKDIR /app
ENV TIME_ZONE Asia/Shanghai
COPY . .
COPY ./api/static .

#RUN apk --update add bind-tools && rm -rf /var/cache/apk/*
RUN go build -v -o /app/skydns ./

ENTRYPOINT ["/app/skydns"]
