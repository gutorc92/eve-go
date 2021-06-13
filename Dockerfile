FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/main

FROM scratch
COPY --from=build /bin/main /bin/main
ENV EVE_GO_URI=""
ENV EVE_GO_DATABASE=""
ENV EVE_GO_FILES=""
EXPOSE 5000
# RUN mkdir /usr/json/
ENTRYPOINT ["/bin/main", "serve"]