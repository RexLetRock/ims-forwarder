BINARY=engine
REPO=registry.hura.asia

FORWARDER=forwarder
FORWARDER_VER=v1.0.0

APP_FORWARDER_IMAGE=${REPO}/${FORWARDER}:${FORWARDER_VER}

run:
	go run service/forwarder/cmd/main.go

build-forwarder:
	GOOS=linux CGO_ENABLED=1 go build -tags musl -o engine service/forwarder/cmd/*.go

docker-build-forwarder:
	CGO_ENABLED=1 docker build -f service/forwarder/zDockerfile -t ${APP_FORWARDER_IMAGE} .
	docker push ${APP_FORWARDER_IMAGE}
	docker image rm ${APP_FORWARDER_IMAGE}
