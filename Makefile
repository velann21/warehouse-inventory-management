run-local-mysql:
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=root -p 3308:3306 -d mysql

run-local-inventory:
	docker run --network=host -d "singaravelan21/inventory_srv:`cat .version`"

build-docker:
	docker build . -t inventory_srv;

push-docker: build-docker
	docker tag "inventory_srv" "singaravelan21/inventory_srv:`cat .version`";
	docker push "singaravelan21/inventory_srv:`cat .version`";



