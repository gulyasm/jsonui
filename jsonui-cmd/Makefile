compilerFlag=-gcflags=-trimpath=$(shell pwd) -asmflags=-trimpath=$(shell pwd)
goFiles=main.go
gomod=off
all: mac linux windows freebsd netbsd openbsd pi

windows:$(goFiles)
	GO111MODULE=$(gomod) GOOS=windows GOARCH=386  go build $(compilerFlag)  -o ../../jsonui-release/win32-jsonui.exe $(goFiles)
	GO111MODULE=$(gomod) GOOS=windows GOARCH=amd64  go build $(compilerFlag)  -o ../../jsonui-release/win64-jsonui.exe $(goFiles)

linux: $(goFiles)
	GO111MODULE=$(gomod) GOOS=linux GOARCH=386  go build $(compilerFlag)  -o ../../jsonui-release/lin32-jsonui $(goFiles)
	GO111MODULE=$(gomod) GOOS=linux GOARCH=amd64  go build $(compilerFlag)  -o ../../jsonui-release/lin64-jsonui $(goFiles)

pi: $(gofiles)
	GO111MODULE=$(gomod) GOOS=linux GOARCH=arm GOARM=6 go build $(compilerFlag)  -o ../../jsonui-release/jsonui-pi $(goFilesU)

mac:
	GO111MODULE=$(gomod) GOOS=darwin GOARCH=amd64  go build $(compilerFlag) -o ../../jsonui-release/mac-jsonui $(goFiles)

freebsd: $(goFiles)
	GO111MODULE=$(gomod) GOOS=freebsd GOARCH=386  go build $(compilerFlag)  -o ../../jsonui-release/freebsd32-jsonui $(goFiles)
	GO111MODULE=$(gomod) GOOS=freebsd GOARCH=amd64  go build $(compilerFlag)  -o ../../jsonui-release/freebsd64-jsonui $(goFiles)

netbsd: $(goFiles)
	GO111MODULE=$(gomod) GOOS=freebsd GOARCH=386  go build $(compilerFlag)  -o ../../jsonui-release/netbsd32-jsonui $(goFiles)
	GO111MODULE=$(gomod) GOOS=freebsd GOARCH=amd64  go build $(compilerFlag)  -o ../../jsonui-release/netbsd64-jsonui $(goFiles)

openbsd: $(goFiles)
	GO111MODULE=$(gomod) GOOS=openbsd GOARCH=386  go build $(compilerFlag)  -o ../../jsonui-release/openbsd32-jsonui $(goFiles)
	GO111MODULE=$(gomod) GOOS=openbsd GOARCH=amd64  go build $(compilerFlag)  -o ../../jsonui-release/openbsd64-jsonui $(goFiles)

solaris: $(goFiles)
	GO111MODULE=$(gomod) GOOS=solaris GOARCH=amd64  go build $(compilerFlag)  -o ../../jsonui-release/sol-jsonui $(goFiles)