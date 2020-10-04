#!/bin/sh

# docker logs -f <container>
# ======== explore this filesystem using bash (for example)
# docker run -t -i yari /bin/bash

dockerTempFolder=~/temp/dockerTemp
dockerFileSource="/data/workspaces/go/src/yari"
executableFile="/data/workspaces/go/bin/yari"

# clean all the current environments
docker container rm yariL1
docker container rm yariF1
docker container rm yariF2
docker container rm yariF3
docker container rm yariF4
# docker container rm yariF5
docker image rm yari

mkdir -p $dockerTempFolder/usr/lib
mkdir -p $dockerTempFolder/usr/bin
cd $dockerTempFolder
cp $dockerFileSource/Dockerfile  ./.
cp $dockerFileSource/config.yaml ./.
cp $executableFile               ./.
# copy libs
cp /usr/lib/libqpid-proton-core.so.10         ./usr/lib/.
cp /usr/lib/x86_64-linux-gnu/libssl.so.1.1    ./usr/lib/.
cp /usr/lib/x86_64-linux-gnu/libsasl2.so.2    ./usr/lib/.
cp /usr/lib/x86_64-linux-gnu/libcrypto.so.1.1 ./usr/lib/.

# #copy db5.3-util from /usr/bin
# cp /usr/bin/db5.3_archive       ./usr/bin/.
# cp /usr/bin/db5.3_checkpoint    ./usr/bin/.
# cp /usr/bin/db5.3_deadlock      ./usr/bin/.
# cp /usr/bin/db5.3_dump          ./usr/bin/.
# cp /usr/bin/db5.3_hotbackup     ./usr/bin/.
# cp /usr/bin/db5.3_load          ./usr/bin/.
# cp /usr/bin/db5.3_log_verify    ./usr/bin/.
# cp /usr/bin/db5.3_printlog      ./usr/bin/.
# cp /usr/bin/db5.3_recover       ./usr/bin/.
# cp /usr/bin/db5.3_replicate     ./usr/bin/.
# cp /usr/bin/db5.3_stat          ./usr/bin/.
# cp /usr/bin/db5.3_upgrade       ./usr/bin/.
# cp /usr/bin/db5.3_verify        ./usr/bin/.
# 
# #copy db-util from /usr/bin
# cp /usr/bin/db_archive          ./usr/bin/.
# cp /usr/bin/db_checkpoint       ./usr/bin/.
# cp /usr/bin/db_deadlock         ./usr/bin/.
# cp /usr/bin/db_dump             ./usr/bin/.
# cp /usr/bin/db_hotbackup        ./usr/bin/.
# cp /usr/bin/db_load             ./usr/bin/.
# cp /usr/bin/db_log_verify       ./usr/bin/.
# cp /usr/bin/db_printlog         ./usr/bin/.
# cp /usr/bin/db_recover          ./usr/bin/.
# cp /usr/bin/db_replicate        ./usr/bin/.
# cp /usr/bin/db_sql              ./usr/bin/.
# cp /usr/bin/db_stat             ./usr/bin/.
# cp /usr/bin/db_upgrade          ./usr/bin/.
# cp /usr/bin/db_verify           ./usr/bin/.
# 
# 
# # copy sasl2-bin from /usr/bin
# cp /etc/default/saslauthd         ./usr/bin/.
# cp /etc/init.d/saslauthd          ./usr/bin/.
# cp /usr/bin/gen-auth              ./usr/bin/.
# cp /usr/bin/sasl-sample-client    ./usr/bin/.
# cp /usr/bin/saslfinger            ./usr/bin/.
# cp /usr/lib/sasl2/berkeley_db.txt ./usr/bin/.
# cp /usr/sbin/sasl-sample-server   ./usr/bin/.
# cp /usr/sbin/saslauthd            ./usr/bin/.
# cp /usr/sbin/sasldbconverter2     ./usr/bin/.
# cp /usr/sbin/sasldblistusers2     ./usr/bin/.
# cp /usr/sbin/saslpasswd2          ./usr/bin/.
# cp /usr/sbin/saslpluginviewer     ./usr/bin/.
# cp /usr/sbin/testsaslauthd        ./usr/bin/.
# 

docker build -t yari ./
docker run -it --net="host" --name=yariL1 yari /YARI/yari -initialRole=1 -nodeId=Leader-1
docker run -it --net="host" --name=yariF1 yari /YARI/yari -initialRole=2 -nodeId=Follower-1
docker run -it --net="host" --name=yariF2 yari /YARI/yari -initialRole=2 -nodeId=Follower-2
docker run -it --net="host" --name=yariF3 yari /YARI/yari -initialRole=2 -nodeId=Follower-3
docker run -it --net="host" --name=yariF4 yari /YARI/yari -initialRole=2 -nodeId=Follower-4
# docker stop yari


# rm -r $dockerTempFolder