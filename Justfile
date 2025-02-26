help:
    @just --list

tidy:
	go mod tidy

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.4 run

sec:
	go run github.com/securego/gosec/v2/cmd/gosec@v2.21.4 ./...

demo:
	go run cmd/example/main.go

clean:
	rm -rf out
