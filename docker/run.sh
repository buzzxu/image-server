#!/bin/bash

CONFIG_FILE=/app/conf.yml
until [ $# -eq 0 ]
do
 case "$1" in
 --domain)
    sed -i "2s/none/$2/g" $CONFIG_FILE
    shift 2;;
 --bodyLimit)
    sed -i "5s/5M/$2/g" $CONFIG_FILE
    shift 2;;
 --sizeLimit)
    sed -i "6s/500K/$2/g" $CONFIG_FILE
    shift 2;;
 --defalut-img)
    sed -i "7s/default.png/$2/g" $CONFIG_FILE
    shift 2;;
 --type)
    sed -i "8s/local/$2/g" $CONFIG_FILE
    shift 2;;
 --maxAge)
    sed -i "4s/31536000/$2/g" $CONFIG_FILE
    shift 2;;
 --jwt-secret)
    sed -i "10s/123456/$2/g" $CONFIG_FILE
    shift 2;;
 --jwt-algorithm)
    sed -i "11s/HS512/$2/g" $CONFIG_FILE
    shift 2;;
 --redis-addr)
    sed -i "18s/host.docker.internal:6379/$2/g" $CONFIG_FILE
    shift 2;;
 --redis-password)
    sed -i "19s/none/$2/g" $CONFIG_FILE
    shift 2;;
 --redis-db)
    sed -i "20s/1/$2/g" $CONFIG_FILE
    shift 2;;
 --redis-pool)
    sed -i "21s/0/$2/g" $CONFIG_FILE
    shift 2;;
 --redis-expire)
    sed -i "22/10800/$2/g" $CONFIG_FILE
    shift 2;;
 --aliyun-endpoint)
    sed -i "24s/http://oss-cn-hangzhou.aliyuncs.com/$2/g" $CONFIG_FILE
    shift 2;;
 --aliyun-accesskey-id)
    sed -i "25s/xux/$2/g" $CONFIG_FILE
    shift 2;;
 --aliyun-accesskey-secret)
    sed -i "26s/xux/$2/g" $CONFIG_FILE
    shift 2;;
 --aliyun-bucket)
    sed -i "27s/xux/$2/g" $CONFIG_FILE
    shift 2;;
 --minio-endpoint)
    sed -i "29s/127.0.0.1:9001/$2/g" $CONFIG_FILE
    shift 2;;
 --minio-accessKey)
    sed -i "30s/xuxiang/$2/g" $CONFIG_FILE
    shift 2;;
 --minio-secretKey)
    sed -i "31s/111111/$2/g" $CONFIG_FILE
    shift 2;;
 --minio-bucket)
    sed -i "32s/buzzxu/$2/g" $CONFIG_FILE
    shift 2;;
 --minio-useSSL)
    sed -i "33s/false/$2/g" $CONFIG_FILE
    shift 2;;
 *) echo " unknow prop $1";shift;;
 esac
done

echo "============config.js==============="
cat $CONFIG_FILE
echo "===================================="

./app