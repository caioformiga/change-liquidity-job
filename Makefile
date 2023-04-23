PKGS ?= $(shell go list ./...)

build-local:
	go build -o app ./main.go

lint:
	$$GOPATH/bin/golint ${PKGS}

test:
	go test -race -covermode=atomic -coverprofile coverage.out -v ${PKGS}

unit-test:	
	go test -p 1 -covermode=atomic -coverprofile coverage.out -v ${PKGS} -tags=unit

run:
	go run main.go

build:
	$$GOPATH/bin/gox -osarch="linux/amd64" -output="app" ./

migrate-create:
	migrate create -ext json -dir migrations -seq ${NAME}

migrate-up:
	migrate -database ${MONGO_URI} -path migrations up

migrate-down:
	migrate -database ${MONGO_URI} -path migrations down

migrate-version:
	migrate -database ${MONGO_URI} -path migrations version

migrate-force:
	migrate -database ${MONGO_URI} -path migrations force ${VERSION}

migrate-create-trading-bots:
	migrate create -ext json -dir migrations_trading_bots -seq ${NAME}

migrate-version-trading-bots:
	migrate -database ${MONGO_URI} -path migrations_trading_bots version

migrate-force-trading-bots:
	migrate -database ${MONGO_URI} -path migrations_trading_bots force ${VERSION}

