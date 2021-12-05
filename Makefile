mod:
	GO111MODULE=on; go mod tidy; go mod vendor;

test:
	go test --coverpkg ./pkg/... -race -covermode atomic -coverprofile=.coverage.out ./pkg/...

cienv:
	ci/actions.sh setup_go

citest:
	ci/actions.sh test_go

cover:
	go tool cover -func=.coverage.out
