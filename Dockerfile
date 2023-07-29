FROM golang:1.17-alpine3.15

LABEL maintainer="@elinana @AntonL @AlpBal @maggieeagle"
WORKDIR /app     

COPY . .
RUN apk add --no-cache gcc musl-dev && go mod download

RUN go build -o forum .
CMD [ "/app/forum" ]