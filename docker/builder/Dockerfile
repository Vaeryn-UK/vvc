FROM golang

RUN apt update

# For hexdump
RUN apt install bsdmainutils

ARG UID

RUN adduser --home /home/dev --uid $UID --disabled-password --gecos "" dev

RUN apt search hexdump

USER dev