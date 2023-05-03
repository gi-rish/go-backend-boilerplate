FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /github.com/NewStreetTechnologies/go-backend-boilerplate

COPY --from=0 github.com/NewStreetTechnologies/go-backend-boilerplate ./

EXPOSE 8210

CMD [ "/github.com/NewStreetTechnologies/go-backend-boilerplate" ]