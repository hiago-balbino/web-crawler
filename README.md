# WEB-CRAWLER
This project was created for learning purposes and is a crawler that goes through the web looking for any information, clicking on every link available. 

Some used tools do not represent the better choice, were just used for learning purposes. For example, MongoDB was used but focusing on performance Redis could be used. The front end was not focused on learning purposes so the `template package` was used.

## Requirements
* [Go](https://golang.google.cn/dl/) 1.19+
* [Docker](https://www.docker.com/products/docker-desktop/)
* [Docker-compose](https://docs.docker.com/compose/install/)
* [GNU Make](https://www.gnu.org/software/make/)
* [golangci-lint](https://golangci-lint.run/)
* [godoc](https://pkg.go.dev/golang.org/x/tools/cmd/godoc)
* [direnv](https://direnv.net/)
    * This is not mandatory but is a easily way to control your environment variable for each project without configuring the variables globally

## Useful commands
You can run the command below to see all the useful commands available and your goals.
```
make help
```
```
help: show this help.
setup: run the command mod download and tidy from Go
vet: run the command vet from Go
tests: run all unit tests
integration-tests: run all integration tests
all-tests: run all unit and integration tests
cover: run the command tool cover to open coverage file as HTML
lint: run all linters configured
sonarqube-up: start sonarqube container
sonarqube-down: stop sonarqube container
sonarqube-analysis: run sonar scanner
fmt: run go formatter recursively on all files
compose-ps: list all containers running
compose-up: start API and dependencies
compose-down: stop API and dependencies
build: create an executable of the application
build-run-api: build project and run the API using the built binary
clean: run the go clean command and removes the application binary
doc: run the project documentation using HTTP
 ```

## How to run de application
To run the project locally you need to export some environment variables and you can do that using `direnv`. You can export the variables below.
```
PORT='8888'
LOG_LEVEL='ERROR'

MONGODB_USERNAME='root'
MONGODB_PASSWORD='example'
MONGODB_DATABASE='crawler'
MONGODB_COLLECTION='page'
MONGODB_PORT='27017'
MONGODB_HOST='mongo'

MONGODB_EXPRESS_USERNAME='root'
MONGODB_EXPRESS_PASSWORD='example'
MONGODB_EXPRESS_PORT='8081'
```

After exporting the environment variables you can run the command `make compose-up`. If you want to run outside of Docker, you can run the command `make build-run-api` and open the address `http://localhost:8888/index`.

If you want to debug the application, you need to export `MONGODB_HOST` variable as `localhost`, comment the `api` service on `docker-compose.yml` and run `make compose-up`. In your IDE you need to set the command as `api` given that the application is using [cobra library](https://github.com/spf13/cobra).

## How to craw page
* Fill the URI and Depth on the form(will be used to limit the depth when fetching pages with so many links that can have lower performance and can take so much time).

## How to run internal documentation
You can do that by running the command `make doc` and accessing the address `http://localhost:6060`.

## How to run sonarqube locally
The project was configured to run `sonarqube` locally and has three commands on Makefile. The `sonarqube` will be downloaded by Docker, but you need to [install the sonar-scanner following your operational system](https://docs.sonarqube.org/latest/analyzing-source-code/scanners/sonarscanner).

To run the `sonarqube` locally you need to export the following environment variables. You can do that using `direnv`.
```
SONAR_PORT='9000'
SONAR_HOST='http://localhost:9000'
SONAR_LOGIN='admin'
SONAR_PASSWORD='admin'
SONAR_BINARY='Here you need to fill it according to your operational system. Example: sonar-scanner for Linux/MacOS or sonar-scanner.bat for Windows'
```

After installing and configuring `sonar-scanner` on `$PATH`(if needed) you will be able to run the commands below. When running the commands `make sonarqube-up` and `make sonarqube-analysis` you can open the address `http://localhost:9000` on your browser and make a login and password as `admin`(maybe `sonarqube` can request to change password).
```
sonarqube-up: start sonarqube container
sonarqube-analysis: run sonar scanner
sonarqube-down: stop sonarqube container
```