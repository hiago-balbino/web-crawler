[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/hiago-balbino/web-crawler.svg)](https://pkg.go.dev/github.com/hiago-balbino/web-crawler)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiago-balbino/web-crawler)](https://goreportcard.com/report/github.com/hiago-balbino/web-crawler)
[![Linter](https://github.com/hiago-balbino/web-crawler/actions/workflows/linter.yml/badge.svg?branch=main&event=push&logo=github&style=flat-square)](https://github.com/hiago-balbino/web-crawler/actions/workflows/linter.yml)
[![Security](https://github.com/hiago-balbino/web-crawler/actions/workflows/security.yml/badge.svg?branch=main&event=push&logo=github&style=flat-square)](https://github.com/hiago-balbino/web-crawler/actions/workflows/security.yml)
[![Tests](https://github.com/hiago-balbino/web-crawler/actions/workflows/tests.yml/badge.svg?branch=main&event=push&logo=github&style=flat-square)](https://github.com/hiago-balbino/web-crawler/actions/workflows/tests.yml)

# 🔍 WEB-CRAWLER
This project was created for learning purposes and is a crawler that go through the web looking for any information by clicking on each available link. 

Some tools used do not represent the best choice, they were only used for learning purposes. For example MongoDB was used, but thinking about performance Redis might be a better alternative. The frontend was not the focus for learning purposes, so the `template package` was used.

## 🧰 Dependencies
* [Go](https://golang.google.cn/dl) 1.19+
* [Docker](https://www.docker.com/products/docker-desktop)
* [Docker-compose](https://docs.docker.com/compose/install)
* [Nginx](https://nginx.org)
* [GNU make](https://www.gnu.org/software/make)
* [Direnv](https://direnv.net)
    * This is not mandatory but is a easily way to control your environment variable for each project without configuring the variables globally
* [Golangci-lint](https://golangci-lint.run)
* [Godoc](https://pkg.go.dev/golang.org/x/tools/cmd/godoc)
* [Govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
* [Viper](https://github.com/spf13/viper)
    * Configuration solution
* [Cobra](https://github.com/spf13/cobra)
    * Library for creating CLI applications
* [Gin web framework](https://github.com/gin-gonic/gin)
* [Gin metrics exporter for Prometheus](https://github.com/penglongli/gin-metrics)
* [Prometheus client golang](https://github.com/prometheus/client_golang)
* [Testify](https://github.com/stretchr/testify)
    * Tools for testifying
* [Testcontainers-Go](https://github.com/testcontainers/testcontainers-go)
* [Zap logger](https://go.uber.org/zap)
* [MongoDB](https://www.mongodb.com)
    * The web-based admin interface Mongo Express was used
* [HTTP expect](https://github.com/gavv/httpexpect)
    * Used to API test
* [Gock HTTP mocking](https://github.com/h2non/gock)
* [Sonarqube](https://www.sonarqube.org)

## 🛠️ Useful commands
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

## ⚙️ Running the Application
To run the project locally you need to export some environment variables and this can be done using `direnv`. You can export the variables below.
```
NGINX_PORT='80'
LOG_LEVEL='ERROR'

MONGODB_USERNAME='root'
MONGODB_PASSWORD='example'
MONGODB_DATABASE='crawler'
MONGODB_COLLECTION='page'
MONGODB_PORT='27017'

MONGODB_EXPRESS_USERNAME='root'
MONGODB_EXPRESS_PASSWORD='example'
MONGODB_EXPRESS_PORT='8081'
```

After exporting the environment variables, you can run the `make compose-up` command and open the `http://localhost/index` address. 

If you want to run the API outside of Docker:
* you need to export the `MONGODB_HOST` variable as `localhost`
* expose the `MONGODB_PORT` in `docker-compose.yml` in the `mongo` service as below
    ```
    ports:
      - ${MONGODB_PORT}:${MONGODB_PORT}
    ```
* comment out the `api` and `nginx` service in `docker-compose.yml` and run `make compose-up` and then `make build-run-api`, the API will run on the default port `http://localhost:8888/index`
* if you want to debug the API, you don't need to run `make build-run-api` and in your IDE you need to set the command to `api` as the application is using [cobra library](https://github.com/spf13/cobra)

## 🏁 How to crawl the page
Fill in the URI and Depth in the form(it will be used to limit the depth when fetching pages with so many links that they can underperform and can take so long).

## 📜 Running Internal Documentation
You can do this by running the `make doc` command and going to the address `http://localhost:6060`.

## 🎯 How to run sonarqube locally
The project is set up to run `sonarqube` locally and has three commands in the Makefile. The `sonarqube` will be downloaded by Docker, but you need to [install sonar-scanner following your operating system](https://docs.sonarqube.org/latest/analyzing-source-code/scanners/sonarscanner).

To run `sonarqube` locally, you need to export the following environment variables. You can do this using `direnv`.
```
SONAR_PORT='9000'
SONAR_HOST='http://localhost:9000'
SONAR_LOGIN='admin'
SONAR_PASSWORD='admin'
SONAR_BINARY='Here you need to fill it according to your operational system. Example: sonar-scanner for Linux/MacOS or sonar-scanner.bat for Windows'
```

After installing and configuring `sonar-scanner` in `$PATH`(if needed) you will be able to run the commands below. By running the `make sonarqube-up` and `make sonarqube-analysis` commands you can open the `http://localhost:9000` address in your browser and login and password as `admin`(perhaps `sonarqube` may prompt you to change your password).
```
sonarqube-up: start sonarqube container
sonarqube-analysis: run sonar scanner
sonarqube-down: stop sonarqube container
```

## 📊 Running the metrics
The project was instrumented using `Prometheus` and `Grafana`, both of which are configured and downloaded through Docker. Prometheus and Grafana will run together with the application, but you need to export the following environment variables below, and you can do this using `direnv`.
```
PROMETHEUS_PORT='9090'
GRAFANA_PORT='3000'
```

The application metrics are exposed using the [ginmetrics library](https://github.com/penglongli/gin-metrics) and can be accessed at `http://localhost:8888/metrics`. These exposed metrics are collected by Prometheus and can be accessed at `http://localhost:9090`. 

The collected metrics are sent to Grafana and can be accessed at `http://localhost:3000`. The default credentials are `admin`/`admin`(Grafana may prompt you to reset the password, but it is optional). After that, you need to configure the `data source` by clicking on the `Configuration` option in the left hand panel and then clicking on `Data source`. Click on the `Add Data Source` button and select `Prometeus` under `Time Series Database`. Fill in the address in the HTTP option as in the image below:

[![datasource](/metrics/docs/images/datasource.png)](/prometheus/docs/images/datasource.png)

After setting up the data source, you can import the file from the dashboard by clicking on the `Dashboard` option in the left side panel and then clicking `+ Import`. You can upload the file placed in this project at `/metrics/grafana/dashboards.json`. After it is loaded, you will see the panels as below:

[![metrics](/metrics/docs/images/metrics.png)](/prometheus/docs/images/metrics.png)