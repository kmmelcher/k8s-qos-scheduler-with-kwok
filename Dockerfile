FROM golang:1.14.4-stretch AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1

WORKDIR /build

# Let's cache modules retrieval - those don't change so often
COPY go.mod .
COPY go.sum .
RUN go mod download -x

# Copy the code necessary to build the application
# You may want to change this to copy only what you actually need.
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build the application
RUN CGO_ENABLED=1 go build -ldflags '-w' -o bin/kube-scheduler cmd/main.go

# Let's create a /dist folder containing just the files necessary for runtime.
# Later, it will be copied as the / (root) of the output image.
WORKDIR /dist
RUN cp /build/bin/kube-scheduler ./kube-scheduler

# Optional: in case your application uses dynamic linking (often the case with CGO),
# this will collect dependent libraries so they're later copied to the final image
# NOTE: make sure you honor the license terms of the libraries you copy and distribute
RUN ldd kube-scheduler | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/

# Create the minimal runtime image
FROM golang:1.14.4-stretch
COPY --from=builder /dist /usr/local/bin/
