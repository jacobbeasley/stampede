# Stage 1: Frontend
FROM node:20-alpine AS frontend
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Backend Builder
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache gcc musl-dev git bash
RUN go install github.com/gobuffalo/cli/cmd/buffalo@latest
ENV GOPROXY=https://proxy.golang.org GO111MODULE=on
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/public/assets ./public/assets
RUN /go/bin/buffalo build --static --skip-assets -o /bin/app

# Stage 3: Release
FROM alpine
RUN apk add --no-cache bash ca-certificates
WORKDIR /bin
COPY --from=builder /bin/app .
COPY database.yml .
ENV ADDR=0.0.0.0
EXPOSE 3000
CMD ["/bin/app"]
