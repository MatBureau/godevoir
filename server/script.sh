docker pull influxdb:3-core

docker run -d --name influxdb3 \
  -p 8181:8181 \
  -v influxdb3-data:/var/lib/influxdb3/data \
  influxdb:3-core
