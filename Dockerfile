FROM golang:alpine
WORKDIR /app
COPY . .
RUN go mod tidy
EXPOSE ${PORT}
CMD [ "go", "run", "." ]
