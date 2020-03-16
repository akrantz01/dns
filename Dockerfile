FROM node:10.16.0-stretch AS frontend-build

WORKDIR /frontend
COPY frontend /

RUN npm install -g yarn
RUN yarn install
RUN yarn build

FROM golang:1.12-stretch AS server-build

RUN go get github.com/GeertJohan/go.rice && go get github.com/GeertJohan/go.rice/rice

RUN mkdir -p src/github.com/akrantz01/krantz.dev/dns/frontend/build
WORKDIR src/github.com/akrantz01/krantz.dev/dns

COPY --from=frontend-build build frontend/build
COPY db ./db
COPY records ./records
COPY roles ./roles
COPY users ./users
COPY util ./util
COPY main.go ./main.go

RUN go get ./...
RUN rice embed-go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /dns -a -installsuffix cgo .

FROM scratch

COPY --from=server-build /dns /dns

ENTRYPOINT ["/dns"]
