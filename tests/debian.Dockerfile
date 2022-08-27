FROM golang:bullseye

EXPOSE 40000

WORKDIR /
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# need to pass -v for a mountpoint when running, mounting the base project
# for example: `docker run --rm -v $PWD:/app envy_ubuntu` will run tests without dlv
# for example: `docker run --rm -v $PWD:/app envy_ubuntu true` will run tests with dlv
