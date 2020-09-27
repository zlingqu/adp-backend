BINARY_NAME=main
srcDir=./src
tag=`date "+%Y-%m-%d-%H-%M"`
adpBackendImage=docker.dm-ai.cn/devops/adp-service:$(tag)

init-env:
	cp -f main $(srcDir)/main.go
	cp -rf $(srcDir)/3rd-api/kubernetes/key $(srcDir) # 各环境kubeconfig目录
	cp -rf $(srcDir)/3rd-api/jenkins/conf $(srcDir)	# jenkins配置文件目录

clean:
	rm -rf build

compile: clean init-env
	cd $(srcDir) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME) -v ./ \
	&& cp -rf key conf ../Makefile build/ \
	&& mv build ../

run:
	cd build && ./$(BINARY_NAME)

docker-build: compile
	docker build -t $(adpBackendImage) -f service-run.Dockerfile .

docker-push:
	docker push $(adpBackendImage)

docker-buildAndTest: docker-build
	docker run -it --rm $(adpBackendImage)

docker-buildAndPush: docker-build
	docker run -it --rm $(adpBackendImage)

docker-buildUseDocker:
	docker build -t $(adpBackendImage) -f service-compile.Dockerfile .



