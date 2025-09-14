FROM node:22.19.0-alpine AS base
LABEL author="Kenley Wang"

RUN apk add --no-cache tzdata
RUN npm install -g pnpm@10.15.1
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"


# Install dependencies
FROM base AS dependencies

WORKDIR /app/web

COPY package.json .
COPY pnpm-lock.yaml .

RUN pnpm install --frozen-lockfile


# build package
FROM base AS builder

WORKDIR /app/web

COPY --from=dependencies /app/web/ .
COPY . .

ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN pnpm build


# run in production
FROM base AS prod

ENV PM2_INSTANCES=2

WORKDIR /app/web

COPY --from=builder /app/web/public ./public
COPY --from=builder /app/web/.next/standalone ./
COPY --from=builder /app/web/.next/static ./.next/static

COPY docker/entrypoint.sh ./entrypoint.sh

RUN pnpm add -g pm2 \
    && mkdir /.pm2 \
    && chown -R 1001:0 /.pm2 /app/web \
    && chmod -R g=u /.pm2 /app/web

EXPOSE 3000
ENTRYPOINT ["/bin/sh", "./entrypoint.sh"]