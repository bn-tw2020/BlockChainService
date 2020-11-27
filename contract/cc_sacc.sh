#!/bin/bash

set -ev

# 설치
docker exec cli peer chaincode install -n sacc -v 1.0.1 -p github.com/sacc1

# 배포
docker exec cli peer chaincode upgrade -n sacc -v 1.0.1 -C mychannel -c '{"Args":["a","100"]}' -P 'OR("Org1MSP.member","Org2MSP.member","Org3MSP.member")'

sleep 3
# 설치 배포 확인
docker exec cli peer chaincode list --instantiated -C mychannel

# 쿼리
docker exec cli peer chaincode query -n sacc -C mychannel -c '{"Args":["get","a"]}'
# 인보크
docker exec cli peer chaincode invoke -n sacc -C mychannel -c '{"Args":["set", "c","300"]}' 
sleep 3
# 쿼리
docker exec cli peer chaincode query -n sacc -C mychannel -c '{"Args":["get","c"]}'
