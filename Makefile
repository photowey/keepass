build:clean
	go env -w GO111MODULE=on
	go env -w GOPROXY=https://goproxy.cn,direct
	go mod download -x

windows:clean
	GOOS=windows GOARCH=amd64 go build

linux:clean
	GOOS=linux GOARCH=amd64 go build

macos:clean
	GOOS=darwin GOARCH=amd64 go build

clean:
	rm -rf keepass

tidy:clean
	go mod tidy -v