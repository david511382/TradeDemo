ENV := local

##########################################
## docker-compose
##########################################

up-test:
	docker-compose up -d

down: # 關閉
	docker-compose down

##########################################
## build
##########################################

build:
	go build .

run:
	CONFIG=test ./zerologix-homework server

##########################################
## test
##########################################

mock:
	mockgen -source="./src/api/interface.go" -destination="./src/api/interface_mock.go" -package=api
	mockgen -source="./src/bll/nats/interface.go" -destination="./src/bll/nats/interface_mock.go" -package=nats
	mockgen -source="./src/bll/trade/interface.go" -destination="./src/bll/trade/interface_mock.go" -package=trade
	mockgen -source="./src/pkg/timeutil/interface.go" -destination="./src/pkg/timeutil/interface_mock.go" -package=timeutil

# -p 設定執行續，預設不同 package 會非同步執行
test: # 測試
	go test \
	./src/pkg/util/... \
	./bootstrap/... \
	./src/repo/... \
	./src/bll/... \
	./src/api \
	./src/server/... \
	--count=1
