FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend

RUN apk add --no-cache \
    autoconf \
    automake \
    libtool \
    nasm \
    make \
    g++ \
    zlib-dev

COPY frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build

FROM golang:1.23.2-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -o trilium-blog .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
ENV TZ="Asia/Shanghai"
WORKDIR /app
COPY --from=backend-builder /app/backend/trilium-blog .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
COPY backend/config.json . 
EXPOSE 8080
CMD ["./trilium-blog"]
