FROM node:16-alpine as build-vue
WORKDIR /app
COPY statics /app

RUN npm install
RUN npm run build

FROM golang:alpine AS builder-app
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go mod download
RUN go build
RUN ls -al /build

#FROM scratch
FROM ubuntu:24.04
ENV TZ=Asia/Seoul \
    DEBIAN_FRONTEND=noninteractive
RUN export PATH=/app:.:$PATH

WORKDIR /app
RUN mkdir -p /app/.aws
RUN /bin/chmod -Rf 400 /app/.aws/
COPY --from=builder-app /build/tz-s3-repo-mgmt /app/
COPY --from=builder-app /build/.env /app/
COPY --from=builder-app /build/config /app/.aws/config
COPY --from=builder-app /build/credentials /app/.aws/credentials
COPY --from=build-vue /app/dist /app/assets
RUN cp -Rf /app/.aws /root/

RUN apt update && \
    apt install curl telnet pkg-config mysql-client dnsutils iputils-ping jq unzip wget vim -y

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-2.0.30.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && ./aws/install --bin-dir /usr/local/bin --install-dir /usr/local/aws-cli --update

RUN curl -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.18.9/2020-11-02/bin/linux/amd64/aws-iam-authenticator
RUN chmod +x aws-iam-authenticator
RUN mv aws-iam-authenticator /usr/local/bin

EXPOSE 8080
ENTRYPOINT ["./tz-s3-repo-mgmt"]
#CMD ["/bin/sh", "-c", "while true; do echo $(date -u) >> out.txt; sleep 5; done"]
