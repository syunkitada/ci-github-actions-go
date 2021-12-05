mod:
	GO111MODULE=on; go mod tidy; go mod vendor;


test:
	go test --coverpkg ./pkg/... -coverprofile=.coverage.out ./pkg/...


cover:
	go tool cover -func=.coverage.out
