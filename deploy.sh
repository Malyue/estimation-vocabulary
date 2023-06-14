#!/bin/bash


echo "Start the deploy,please wait a moment......"

imageName="estimation-vocabulary"
containerName="estimation-vocabulary"

ARG1=$(docker ps -aqf name=${containerName})
ARG2=$(docker images -q --filter reference=${imageName})

## if it is not null,stop and delete the container
if [  -n "$ARG1"  ]; then
  docker rm -f $(docker stop $ARG1)
  echo "$containerName container deleted"
fi

if [  -n "$ARG2"  ]; then
  docker rmi -f $ARG2
  echo "$imageName image deleted"
fi

docker rmi $(docker images | grep "none" |awk '{print $3}')
echo "image deleted"

docker build -t ${imageName} .

# 具体端口号可以自己指定
docker run --name ${containerName} -d -p 8080:8080 ${imageName}

echo "build success"
