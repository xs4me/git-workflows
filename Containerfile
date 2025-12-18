ARG VERSION=latest

FROM golang:1.25-alpine AS build
ARG VERSION

WORKDIR /app/

COPY src/go.mod src/go.sum ./
RUN go mod verify && go mod download
COPY src/ .

RUN GOOOS=$TARGETOS GOARCH=$TARGETARCH go build -o git-workflows -ldflags="-X main.version=$VERSION" .

FROM alpine:3.23.2

RUN apk add --no-cache  \
    bash  \
    ca-certificates  \
    curl  \
    git  \
    openssh \
    wget

RUN addgroup -g 1000 -S workflow && \
    adduser -u 1000 -S workflow -G workflow

RUN mkdir -p /workflow/

COPY --from=build /app/git-workflows /bin/git-workflows
COPY src/templates/default-descriptor.json /workflow/default-descriptor.json

RUN chgrp -R 0 /workflow  && \
    chmod -R g=u /workflow/ && \
    chgrp -R 0 /bin/git-workflows &&  \
    chmod -R g=u /bin/git-workflows

USER 1000

ENTRYPOINT [ "/bin/git-workflows" ]