SHELL = /bin/bash
export GO111MODULE=on

NAME = geoipserver
OUTPUT_DIR = bin

.PHONY: clean server

clean:
	@rm -rf ${OUTPUT_DIR}

server: clean
	@mkdir -p ${OUTPUT_DIR}
	go build -v -ldflags="-s -w" -o ${OUTPUT_DIR}/${NAME} ./cmd/${NAME}/...

docker-image: clean
	docker build -t ${NAME} -f Dockerfile .
