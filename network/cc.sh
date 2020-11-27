#!/bin/bash

# 체인코드 설치
docker exec cli peer chaincode install -n fabcar -v 1.1 -p github.com/fabcar/go

# 체인코드 배포
docker exec cli peer chaincode upgrade -n fabcar -v 1.1 -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member")'

sleep 3

# 체인코드 설치 배포 확인
docker exec cli peer chaincode list --installed
docker exec cli peer chaincode list --instantiated -C mychannel

# 체인코드 테스트 initLedger
docker exec cli peer chaincode invoke -n fabcar -C mychannel -c '{"Args":["initLedger"]}'
sleep 3
# 체인코드 테스트 queryAllCars
docker exec cli peer chaincode query -n fabcar -C mychannel -c '{"Args":["queryAllCars"]}'
# 체인코드 테스트 createCar
docker exec cli peer chaincode invoke -n fabcar -C mychannel -c '{"Args":["createCar", "CAR10", "bmw","320d","white","ckh"]}'
sleep 3
# 체인코드 테스트 queryCar
docker exec cli peer chaincode query -n fabcar -C mychannel -c '{"Args":["queryCar","CAR10"]}'
