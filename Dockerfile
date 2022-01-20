FROM golang:alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . /app

# Build
RUN go build  -o . ./...

# Run
CMD [ "./download2json" ]