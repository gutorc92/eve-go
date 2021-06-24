FROM public.ecr.aws/lambda/provided:al2 as build-lambda
# install compiler
RUN yum install -y golang
RUN go env -w GOPROXY=direct
# cache dependencies
ADD go.mod go.sum ./
RUN go mod download
# build
ADD . .
RUN go build -o /main
# copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2 as lambda
COPY --from=build-lambda /main /main
EXPOSE 8080
# COPY ./json/farm.json ./farm.json
ENTRYPOINT [ "/main", "lambda" ]

FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY . /src/
RUN mkdir /usr/json/
RUN CGO_ENABLED=0 go build -o /bin/main

FROM alpine
COPY --from=build /bin/main /bin/main
COPY --from=build /usr/json /usr/json
# COPY json/module.json /usr/json
ENV EVE_GO_URI=""
ENV EVE_GO_DATABASE=""
ENV EVE_GO_FILES="/usr/json"
EXPOSE 5000
ENTRYPOINT ["/bin/main", "serve"]         