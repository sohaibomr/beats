FROM phusion/baseimage:latest

RUN apt-get update
RUN apt-get -y -q install libpcap0.8 wget

RUN mkdir /packetbeat

COPY packetbeat /packetbeat/
COPY ./packetBeatRunner/pbrunner /packetbeat/
COPY configs /packetbeat/configs
COPY configs/entrypoint.sh /packetbeat/entrypoint.sh
COPY configs/packetbeat.yml /packetbeat/packetbeat.yml
RUN chmod go-w /packetbeat/packetbeat.yml
WORKDIR /packetbeat
RUN ls

CMD ["bash", "entrypoint.sh"]
