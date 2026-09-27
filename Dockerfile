FROM node:24.21.0-bookworm-slim AS frontend
WORKDIR /src
COPY package.json package-lock.json .npmrc ./
RUN npm ci
COPY tsconfig.json vite.config.ts ./
COPY web ./web
RUN npm run build

FROM golang:1.27.1-bookworm AS backend
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -o /out/server ./cmd/server && \
    CGO_ENABLED=0 go build -trimpath -o /out/manage ./cmd/manage

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=backend /out/server /out/manage /app/
COPY --from=frontend /src/web/build /app/web/build
COPY web/public/fonts /app/web/public/fonts
ENV APP_ENV=production HTTP_ADDR=0.0.0.0:3000
EXPOSE 3000
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
