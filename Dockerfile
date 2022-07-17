FROM golang:1.18

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# TODO Only copy frontend/build instead
COPY ./scripts ./scripts
COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY ./frontend/build ./frontend/build

RUN ./scripts/build.sh

COPY accounts.json ./

CMD ["./build/server"]
