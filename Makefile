# note: call scripts from /scripts

image_name =go-tech
harbor_addr=harbor.apulis.cn:8443/aistudio/app/${image_name}
tag        =$(tag)
arch       =$(shell arch)

testarch:
ifeq (${arch}, x86_64)
	@echo "current build host is amd64 ..."
	$(eval arch=amd64)
else ifeq (${arch},aarch64)
	@echo "current build host is arm64 ..."
	$(eval arch=arm64)
else
	echo "cannot judge host arch:${arch}"
	exit -1
endif
	@echo "arch type:$(arch)"





get-deps:
	#git submodule sync
	#git submodule update --init --recursive
	export GOPRIVATE="gitlab.apulis.com.cn"
	export GOPROXY="https://goproxy.cn"
	export  GO111MODULE=on

	export  GITLAB_USER="linjian"
	export  GITLAB_PWD="Lj3834350"

	git config url."https://${GITLAB_USER}:${GITLAB_PWD}@${GOPRIVATE}".insteadOf "https://${GOPRIVATE}"
	GOPRIVATE="gitlab.apulis.com.cn" go mod tidy
	go mod download

vet-check-all: get-deps
	go vet ./...

gosec-check-all: get-deps
	gosec ./...

bin: get-deps
	go build -o ${image_name} cmd/${image_name}.go

run: bin
	./go-tech

docker:
	docker build -f Dockerfile . -t ${image_name}:v0.1.0

gen-swagger:
	swag init -g cmd/${image_name}.go -o api

builder:
	docker build -t ${image_name} -f build/Dockerfile .

push:
	docker tag ${image_name} ${harbor_addr}:${tag}
	docker push ${harbor_addr}:${tag}

dist: testarch
	docker build -t ${image_name} -f build/Dockerfile .
	docker tag ${image_name} ${harbor_addr}/${arch}:${tag}
	docker push ${harbor_addr}/${arch}:${tag}

manifest:
	./docker_manifest.sh ${harbor_addr}:${tag}

localpush:
	docker build -t ${image_name} -f buildlocal/Dockerfile .
	docker tag ${image_name} ${harbor_addr}:${tag}
	docker push ${harbor_addr}:${tag}

localbuild: get-deps bin

dkpush: localbuild localpush

gitlab:
	git submodule update --init
	docker build -t ${image_name} -f build/Dockerfile .
	docker tag ${image_name} ${harbor_addr}:${tag}
	docker push ${harbor_addr}:${tag}

cicdbuild:
	docker build --add-host gitlab.apulis.com.cn:120.79.42.225 --build-arg GITLAB_USER=${GITLAB_USER} --build-arg GITLAB_PWD=${GITLAB_PWD} -t ${image_name} -f build/Dockerfile .
	docker tag ${image_name} ${harbor_addr}:${tag}
	docker push ${harbor_addr}:${tag}
