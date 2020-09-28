run-local-mysql:
	./scripts.sh

change-permission-scripts.sh:
	chmod 700 ./scripts.sh

build-docker:
	docker build . -t inventory_srv_latest

