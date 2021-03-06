FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY main.go handlers.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/main

FROM scratch
COPY --from=build /bin/main /bin/main
ENTRYPOINT ["/bin/main"]