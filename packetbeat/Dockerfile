FROM phusion/baseimage:latest

RUN apt-get update
RUN apt-get -y -q install libpcap0.8 wget

RUN mkdir /packetbeat

COPY packetbeat /packetbeat/
COPY ./packetBeatRunner/pbrunner /packetbeat/
COPY configs /packetbeat/configs
WORKDIR /packetbeat
RUN ls

CMD ["./pbrunner"]