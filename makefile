PROJECT_NAME := stress-test
PROJECT_IMAGE := $(PROJECT_NAME):latest
BIN_NAME := $(shell echo $(PROJECT_NAME) | sed -r 's/-/_/g')
TEST_URL =https://quotes.toscrape.com/,https://crawler-test.com,https://toscrape.com/,https://www.scrapethissite.com/,https://crawler-test.com/,https://httpbin.org/,https://the-internet.herokuapp.com/,http://books.toscrape.com/,https://realpython.github.io/fake-jobs/,https://s1.demo.opensourcecms.com/wordpress/

DOCKER_RUN := docker run --network=host $(PROJECT_IMAGE)

url-%:
	@echo $(shell echo $(TEST_URL) | cut -d',' -f$*)

run-url-%: build-docker
	@$(eval requests ?= 1000)
	@$(eval concurrency ?= 100)
	$(DOCKER_RUN) -url=$(shell echo $(TEST_URL) | cut -d',' -f$*) -requests=$(requests) -concurrency=$(concurrency)

build-docker:
	docker build --target production --tag $(PROJECT_IMAGE) .

build-go:
	go build -o $(BIN_NAME) main.go

docker-run: build-docker
	$(DOCKER_RUN) -url=$(url) -requests=$(requests) -concurrency=$(concurrency)

go-run: build-go
	./$(BIN_NAME) -url=$(url) -requests=$(requests) -concurrency=$(concurrency)