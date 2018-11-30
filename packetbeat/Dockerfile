FROM phusion/baseimage:latest

RUN apt-get update
RUN apt-get -y -q install libpcap0.8 wget

RUN mkdir /packetbeat

COPY packetbeat /packetbeat/
COPY ./packetBeatRunner/pbrunner /packetbeat/
COPY configs /packetbeat/configs
WORKDIR /packetbeat
RUN ls

# CMD ["./packetbeat","-e", "-E","output.elasticsearch.hosts=['http://${ES_SERVER}:9200']", "-E", "packetbeat.interfaces.device=${INTERFACE}", "-c", "configs/packetbeat.yml"]
CMD ["./pbrunner"]