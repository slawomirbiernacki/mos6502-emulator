FROM debian:11

RUN echo "deb     http://debian.trikaliotis.net/ stable contrib" >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get install cc65

WORKDIR assembler/ca65