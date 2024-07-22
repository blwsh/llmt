FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS build
WORKDIR /src
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/llmt github.com/blwsh/llmt/cmd/llmt

FROM --platform=$BUILDPLATFORM golang:1.22-alpine
COPY --from=build /out/llmt /bin
ENTRYPOINT ["/bin/llmt"]
