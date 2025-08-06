# ---- Build stage ----
FROM node:22.18.0-alpine AS builder
WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run server:build

# ---- Runtime stage ----
FROM node:22.18.0-alpine

WORKDIR /app
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist

RUN chown -R node:node /app

USER node

EXPOSE 3000
CMD ["node", "dist/src/index.js"]