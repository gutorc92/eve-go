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