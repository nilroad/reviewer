FROM basereg.oceantim.com/base/devops/golang/golang-base-images/builder:1.0.8 AS build
WORKDIR /app
COPY . .
RUN  --mount=type=secret,id=netrc,dst=/srv/app/.netrc \
    --mount=type=cache,target=/root/.cache/go-build \
    set -eux; \
    cp /srv/app/.netrc ~/.netrc && \
    go env -w GOPROXY=https://package.oceantim.com/repository/go-golang.org,https://proxy.golang.org,direct && \
    go env -w GOPRIVATE="git.oceantim.com/*" && \
    go generate ./... && \
    CGO_ENABLED=0 go build -trimpath -v -a -ldflags "-w -s" -o build/ ./cmd/...

FROM basereg.oceantim.com/base/devops/golang/golang-base-images/runner:1.0.8
ENV TZ=Asia/Tehran
WORKDIR /app/
COPY --from=build /app/build/ .
COPY ./deploy/ ./deploy/
COPY ./migrations/ ./migrations/
COPY ./deploy/docker/general/configs/supervisord.conf /etc/supervisor/supervisord.conf
ENTRYPOINT ["/app/deploy/docker/general/scripts/entrypoint.sh"]
