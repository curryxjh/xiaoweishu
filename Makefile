.PHONY: docker
docker:
	@rm webook || ture
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f xjh/webook:v0.0.1
	@docker build -t xjh/webook:v0.0.1 .