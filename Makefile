BINARY=engine
REPO=registry.hura.asia
REPO_DEV=14.225.17.203:5000

FORWARDER=forwarder
FORWARDER_VER=v1.0.1

APP_FORWARDER_IMAGE=${REPO}/${FORWARDER}:${FORWARDER_VER}
APP_FORWARDER_IMAGE_DEV=${REPO_DEV}/${FORWARDER}:${FORWARDER_VER}

run:
	go run service/forwarder/cmd/main.go

build-forwarder:
	GOOS=linux CGO_ENABLED=1 go build -tags musl -o engine service/forwarder/cmd/*.go

docker-run:
	CGO_ENABLED=1 docker build -f service/forwarder/zDockerfile -t ${APP_FORWARDER_IMAGE}_local .
	docker run -d -p 19000:19000 ${APP_FORWARDER_IMAGE}_local

docker-build-forwarder:
	CGO_ENABLED=1 docker build -f service/forwarder/zDockerfile -t ${APP_FORWARDER_IMAGE} .
	docker push ${APP_FORWARDER_IMAGE}
	docker image rm ${APP_FORWARDER_IMAGE}

docker-build-forwarder-dev:
	CGO_ENABLED=1 docker build -f service/forwarder/zDockerfile -t ${APP_FORWARDER_IMAGE_DEV} .
	docker push ${APP_FORWARDER_IMAGE_DEV}
	docker image rm ${APP_FORWARDER_IMAGE_DEV}

docker-run-forwarder-dev:
	docker run -d -p 19000:19000 ${APP_FORWARDER_IMAGE_DEV}