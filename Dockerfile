FROM debian:bookworm

RUN apt-get update
RUN apt-get install -y webp ffmpeg

ENV SKIP_DOWNLOAD=true
ENV VENDOR_PATH=/usr/bin

COPY cdn /usr/bin/cdn

ENTRYPOINT ["/usr/bin/cdn"]