run-server:
	go run main.go

run-frontend:
	yarn start

node_modules:
	yarn install

build-frontend:
	yarn build

build-docker:
	docker build -t whatthecard .

build-docker-dev:
	docker build -t whatthecard:dev -f Dockerfile.dev .
