FROM mirror-pub.camera360.com/base/golang-builder:1.17.0 as builder
COPY . /app
WORKDIR /app
ENV CGO_ENABLED=0
RUN /bin/sh -c 'go mod tidy && make build'

# 运维使用的分割线
#---DoNotDelete

#FROM alpine:3.13
FROM mirror-pub.camera360.com/base/centos7.8:basic
WORKDIR /app
COPY --from=builder /app/bin/app /app/bin/app
ADD configs /app/configs
EXPOSE 8000/tcp
#ENTRYPOINT [ "/app/bin/app" ]
