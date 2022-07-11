build:
	hack/verify.sh
	go build -o bin/rpi-humidity cmd/rpi-humidity/rpi-humidity.go

all: build

clean:
	rm -r bin
