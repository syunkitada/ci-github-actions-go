mod:
	GO111MODULE=on; go mod tidy; go mod vendor;

test:
	go test --coverpkg ./pkg/... -coverprofile=.coverage.out ./pkg/...

cienv:
	ci/env.sh setup_go

citest:
	. /opt/ci/envgorc && cd ${MODULE_PATH} && go test --coverpkg ./pkg/... -coverprofile=.coverage.out ./pkg/...
	. /opt/ci/envgorc && cd ${MODULE_PATH} && gcov2lcov-linux-amd64 -infile .coverage.out -outfile /tmp/coverage.lcov

cover:
	go tool cover -func=.coverage.out
