FROM node:24.14.1-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci
COPY frontend .
RUN npm run build

FROM golang:1.25.9-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -o trilium-blog .

FROM alpine:3.21.3
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/backend/trilium-blog .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
VOLUME /app/custom /app/data
EXPOSE 8080
CMD ["./trilium-blog"]
