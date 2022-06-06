# to build this docker image:
#   docker build . -t ghcr.io/jonathongardner/gocv:4.5.5-0.0.0
FROM gocv/opencv:4.5.5 as build

ENV GOPATH /go

COPY dynamic-link-tar.go /go/src/dynamic-link-tar/
WORKDIR /go/src/dynamic-link-tar
RUN GO111MODULE=off go build -o /bin/dynamic-link-tar

COPY mjpeg-streamer/ /go/src/mjpeg-streamer/

WORKDIR /go/src/mjpeg-streamer
RUN go build

RUN mkdir /build && dynamic-link-tar mjpeg-streamer out.tar && tar -xf out.tar -C /build

FROM scratch
COPY --from=build /build /
ENV LD_LIBRARY_PATH=/usr/local/lib/

ENTRYPOINT ["/mjpeg-streamer"]
