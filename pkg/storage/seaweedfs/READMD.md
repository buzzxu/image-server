#### install seaweedfs
```shell script
# 使用filer 依赖levelDB
# 需要安装cmake
yum install gcc-c++ openssl-devel
./bootstrap
make
sudo make install
# 安装snappy levelDB 

wget https://github.com/chrislusf/seaweedfs/releases/download/1.76/linux_amd64.tar.gz
mkdir -p "/data/files/data"
mkdir -p "/data/files/vol/vol1" && chmod 755 /data/files/vol/vol1
mkdir -p "/data/files/vol/vol2" && chmod 755 /data/files/vol/vol2
nohup ./weed master -mdir=/data/files/data -port=9333 -defaultReplication="001" -whiteList="" > /data/logs/weed/master.log 2>&1 &
nohup ./weed volume -dataCenter=dc1 -dir="/data/files/vol/vol1" -max=5 -mserver="192.168.0.135:9333" -port=8080 > /data/logs/weed/vol1.log 2>&1 &
nohup ./weed volume -dataCenter=dc1 -dir="/data/files/vol/vol2" -max=5 -mserver="192.168.0.135:9333" -port=8081 > /data/logs/weed/vol2.log 2>&1 &
./weed scaffold -config=filer -output="."
nohup ./weed filer -master="192.168.0.135:9333" -port=18880 -defaultReplicaPlacement='001' > /data/logs/weed/filer.log 2>&1 &
yum install -y fuse
nohup ./weed mount -filer=192.168.0.135:18880 -dir=/data/files/images > /data/logs/weed/mount.log 2>&1 &
```

#### local debuger

```shell script
brew cask install osxfuse
brew install sshfs
sudo chown ${whoami}  /usr/local/share/man/man1 
brew link sshfs
```
```shell script
chmod a+x data/weed/vol1
chmod a+x data/weed/vol2

nohup ./bin/weed master -mdir=data/weed/data -port=9333 -defaultReplication="001"  > data/logs/weed/master.log 2>&1 &
nohup ./bin/weed volume -dataCenter=dc1 -dir="data/weed/vol1" -max=5 -mserver="127.0.0.1:9333" -port=18080 > data/logs/weed/vol1.log 2>&1 &
nohup ./bin/weed volume -dataCenter=dc1 -dir="data/weed/vol2" -max=5 -mserver="127.0.0.1:9333" -port=18081 > data/logs/weed/vol2.log 2>&1 &
nohup ./bin/weed filer -master=127.0.0.1:9333 -port=18880 > data/logs/weed/filer.log 2>&1 &
# 挂载
./bin/weed mount -filer=127.0.0.1:18880 -dir=/Users/xux/data/images &
```