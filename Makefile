run-local-mysql:
	./scripts.sh

change-permission-scripts.sh:
	chmod 700 ./scripts.sh

build-docker:
	docker build . -t inventory_srv;
push-docker: build-docker
	docker tag "inventory_srv" "singaravelan21/inventory_srv:`cat .version`";
	docker push "singaravelan21/inventory_srv:`cat .version`";



