FROM alpine:3.20

WORKDIR /usr/src/app
RUN apk add --no-cache go

# Dependencies
COPY client/go.mod .
COPY client/go.sum .
RUN go mod download

# Add source
COPY client .
RUN go build -o /usr/local/bin/client .

# Run the Go program
CMD [ "/usr/local/bin/client" ]
