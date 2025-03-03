FROM debian:bookworm

RUN apt-get update
RUN apt-get install -y webp ffmpeg

ENV SKIP_DOWNLOAD=true
ENV VENDOR_PATH=/usr/bin

COPY amur /usr/bin/amur

ENTRYPOINT ["/usr/bin/amur"]