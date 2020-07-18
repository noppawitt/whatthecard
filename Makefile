run-server:
	go run main.go

run-frontend:
	yarn --cwd ./web serve

node_modules:
	yarn --cwd ./web install

build-frontend:
	yarn --cwd ./web build

build-docker:
	docker build -t whatthecard .

build-docker-dev:
	docker build -t whatthecard:dev -f Dockerfile.dev .
