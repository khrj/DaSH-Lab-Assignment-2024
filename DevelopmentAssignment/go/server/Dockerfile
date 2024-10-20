FROM alpine:3.20

WORKDIR /usr/src/app
RUN apk add --no-cache go

# Dependencies
COPY server/go.mod .
COPY server/go.sum .
RUN go mod download

# Add source
COPY server .
RUN go build -o /usr/local/bin/server .

# Run the Go program
CMD [ "/usr/local/bin/server" ]
