APPNAME=phd

build:
	GO111MODULE=on go build -o bin/${APPNAME} main.go

clean:
	rm bin/${APPNAME}* bin/*.yml