#### Quickstart
```shell script
docker run -p 9000:9000 --name minio \
  -v /Users/xux/data/minio:/data \
  minio/minio server /data
```
#### DEPLOYMENT Minio

```shell script
cd /data/files
mkdir data1-1  data1-2  data2-1  data2-2  data3-1  data3-2  data4-1  data4-2

cd /data/docker_vm/minio
docker-compose up -d
curl http://127.0.0.1:9001
```