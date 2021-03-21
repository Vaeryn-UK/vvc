FROM vvc-builder AS builder

COPY ./ /vvc

WORKDIR /vvc

RUN go install ./cmd/...

FROM debian

COPY --from=builder /go/bin/compiler /bin/compile
COPY --from=builder /go/bin/execute /bin/execute

COPY ./programs /programs