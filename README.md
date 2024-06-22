# Service

## Design Considerations

- `Mysql`
  - as we are storing service's name, version and description only, it would be easy to maintain in relational database.
- `GIN`
  - performant, we renowned in opensource community, easy to use and battle tested at my end. 
- `testcontainers-go`
  - `Testcontainers-Go` automates container lifecycle management for integration testing, enabling effortless creation of isolated environments for database testing without manual setup or teardown.
  - It ensures consistency across testing environments, enhances test reproducibility, and enables efficient debugging by closely mirroring production setups
- `GORM`
  - GORM simplifies database operations by allowing developers to interact with the database using Go structures and methods rather than writing raw SQL queries.
  - GORM supports soft deletes, where records are not actually removed from the database but are marked as deleted, allowing for easy recovery and audit.


### Data Model
| Field | Type | Null | Key | Default | Extra |
| :--- | :--- | :--- | :--- |:--------| :--- |
| id | bigint unsigned | NO | PRI | null    | auto\_increment |
| created\_at | datetime\(3\) | YES |  | null    |  |
| updated\_at | datetime\(3\) | YES |  | null    |  |
| deleted\_at | datetime\(3\) | YES |  | null    |  |
| name | varchar\(50\) | YES |  | null    |  |
| description | longtext | YES |  | null    |  |
| version | int | YES |  | 1       |  |
| is\_active | tinyint\(1\) | YES |  | 1       |  |
| tags | varchar\(255\) | YES |  | null    |  |




## Assumptions
- The service `Name` is not unique.
- `Version` is an `int`, e.g., `1`, `9`.
- The service version is server-side driven.
- The first version of the service will always be `1`.
- The `Create service` call will be used to create the first version. For creating a new version, use `PATCH /api/v1/services/{name}/{version}`.
- We use soft deletion in the DB.
- We provided a config file and environment variable support so that it can be hosted on both containerized and non-containerized environments.
- The service name's maximum length is 50 characters.
- We haven't put any limit on service versions.

## Essential Add-ons for Production Readiness.
- File-based MySQL migrations.
- Added `,` separated tags in the model to fine-tune searching.
- Integration with `swaggo/swag` to generate automated swagger documentation from comments.
- Added GIN recovery middleware to recover from unintended panic errors.
- Integrated `Viper` for Config support. (file/env vars)
- Implemented service draining period for graceful shutdowns, the default is `30` seconds. (`The intention is to tackle spot interruptions in the cloud.`)
- Health check points on `/healthcheck` along with `/liveness` and `/readiness`.
- GORM allows developers to interact with the database using Go structures and methods rather than writing raw SQL queries.
- Dockerfile
  - Dependency Caching
  - Multi-stage Build
  - Minimal Base Image (`gcr.io/distroless/static-debian11`)
  - Compiling the Go binary statically (CGO_ENABLED=0) ensures that the application does not depend on any shared libraries, making it more portable and easier to run in different environments.
  
## Integration Tests
Currently, there is one test in `main_test.go`. This file can be extended for more test cases.
`TestShouldCheckGetAllServiceResponse`
```
❯ go clean -i -cache && go test -v -tags='!exclude_tests' ./...                                                                                                                                                                      ─╯
?       github.com/suyog1pathak/services/api/v1/generic [no test files]
?       github.com/suyog1pathak/services/api/v1/healthcheck     [no test files]
?       github.com/suyog1pathak/services/api/v1/response        [no test files]
?       github.com/suyog1pathak/services/docs   [no test files]
?       github.com/suyog1pathak/services/internal/healthcheck   [no test files]
?       github.com/suyog1pathak/services/internal/service       [no test files]
?       github.com/suyog1pathak/services/migration      [no test files]
?       github.com/suyog1pathak/services/pkg/config     [no test files]
?       github.com/suyog1pathak/services/pkg/controllers        [no test files]
?       github.com/suyog1pathak/services/pkg/datastore  [no test files]
?       github.com/suyog1pathak/services/pkg/errors/service     [no test files]
?       github.com/suyog1pathak/services/pkg/logger     [no test files]
?       github.com/suyog1pathak/services/pkg/middleware/healthcheck     [no test files]
?       github.com/suyog1pathak/services/pkg/middleware/service [no test files]
?       github.com/suyog1pathak/services/pkg/model      [no test files]
?       github.com/suyog1pathak/services/pkg/server     [no test files]
?       github.com/suyog1pathak/services/pkg/util       [no test files]
2024/06/11 23:54:15 Setting up test infra..
2024/06/11 23:54:15 github.com/testcontainers/testcontainers-go - Connected to docker: 
  Server Version: 24.0.5
  API Version: 1.43
  Operating System: Docker Desktop
  Total Memory: 9699 MB
  Resolved Docker Host: unix:///Users/suyog.pathak/.docker/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: f4b96813db1578e4c6aa2a0bf3fb84b46f5ebcccb9de15c9caad6d71dee5b784
  Test ProcessID: 2b041411-ad29-401b-8ed5-be4f8d37f8e7
2024/06/11 23:54:15 🐳 Creating container for image testcontainers/ryuk:0.7.0
2024/06/11 23:54:15 ✅ Container created: 885b6b34e5be
2024/06/11 23:54:15 🐳 Starting container: 885b6b34e5be
2024/06/11 23:54:15 ✅ Container started: 885b6b34e5be
2024/06/11 23:54:15 🚧 Waiting for container id 885b6b34e5be image: testcontainers/ryuk:0.7.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms}
2024/06/11 23:54:15 🔔 Container is ready: 885b6b34e5be
2024/06/11 23:54:15 🐳 Creating container for image mysql:8.0.32
2024/06/11 23:54:15 ✅ Container created: d407b2cb536a
2024/06/11 23:54:15 🐳 Starting container: d407b2cb536a
2024/06/11 23:54:15 ✅ Container started: d407b2cb536a
2024/06/11 23:54:15 🚧 Waiting for container id d407b2cb536a image: mysql:8.0.32. Waiting for: &{timeout:<nil> Log:port: 3306  MySQL Community Server IsRegexp:false Occurrence:1 PollInterval:100ms}
2024/06/11 23:54:22 🔔 Container is ready: d407b2cb536a
2024/06/11 23:54:22 🐳 Starting container: d407b2cb536a
2024/06/11 23:54:22 ✅ Container started: d407b2cb536a
2024/06/11 23:54:22 🚧 Waiting for container id d407b2cb536a image: mysql:8.0.32. Waiting for: &{timeout:<nil> Log:port: 3306  MySQL Community Server IsRegexp:false Occurrence:1 PollInterval:100ms}
2024/06/11 23:54:22 🔔 Container is ready: d407b2cb536a
DEBUG------> 172.17.0.3   63612
2024/06/11 23:54:22 running migrations.
Running db migrations..
=== RUN   TestShouldCheckGetAllServiceResponse
{"time":"2024-06-11T23:54:22.364142+05:30","level":"INFO","msg":"please access swagger docs","path":"http://localhost:8080/docs/index.html"}
{"time":"2024-06-11T23:54:22.364394+05:30","level":"INFO","msg":"received a request to get all services."}
{"time":"2024-06-11T23:54:22.36441+05:30","level":"DEBUG","msg":"fetching all services"}
{"time":"2024-06-11T23:54:22.36856+05:30","level":"DEBUG","msg":"no error at middleware"}
{"time":"2024-06-11T23:54:22.368614+05:30","level":"INFO","msg":"Incoming request","request":{"time":"2024-06-11T23:54:22.364235+05:30","method":"GET","host":"","path":"/api/v1/services","query":"","params":{},"route":"/api/v1/services","ip":"","referer":"","length":0},"response":{"time":"2024-06-11T23:54:22.368583+05:30","latency":4348500,"status":200,"length":2},"id":"7426eb13-b967-4767-956a-08edce69cda7"}
--- PASS: TestShouldCheckGetAllServiceResponse (0.01s)
PASS
2024/06/11 23:54:22 Tearing down test infra..
2024/06/11 23:54:22 🐳 Terminating container: d407b2cb536a
2024/06/11 23:54:22 🚫 Container terminated: d407b2cb536a
ok      github.com/suyog1pathak/services        7.616s
```

## Trade-offs
- haven't put limit on versions.
- Table level indexing is yet to be implemented. 

## Pending / future scope.
- authentication/authorization on the API

## Usage
### How to run locally
- setup local mysql with database
```
docker run -it -d  -p "0.0.0.0:3306:3306" -e MYSQL_ROOT_PASSWORD="test123" mysql:8.0.32
 
❯ mysql -h 127.0.0.1 -u root -p                                                                                                                          
Enter password: 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 8
Server version: 8.0.32 MySQL Community Server - GPL

Copyright (c) 2000, 2024, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> create database services;
Query OK, 1 row affected (0.00 sec)

```

- update the `config/config.yaml` with creds or set env vars with prefix `APP`

e.g

```
db:
  user: "test"
  password: "password1"
  host: "192.168.0.5"
  name: "services"
  port: 3306
```
```
export APP_APP_HTTP_PORT=8080
export APP_DB_HOST=0.0.0.0
export APP_DB_NAME=services
export APP_DB_PASSWORD=test123
export APP_DB_USER=root
```

- run db migrations
```
> go run cmd/run-db-migrations.go                                                                                         ─╯
Running db migrations..
```

- start server
```
go run cmd/run-services.go 
```

- (optional) start service as docker container
```
❯ docker build -t services:latest .
...... output truncated ......


❯ docker run -p 8080:8080 -it -e APP_DB_HOST=192.168.29.177 -e APP_DB_USER=root -e APP_DB_NAME=services -e APP_DB_PASSWORD=test123 -e GIN_MODE=release services:latest              
{"time":"2024-06-11T16:07:51.119074554Z","level":"INFO","msg":"starting server at","port":8080}
{"time":"2024-06-11T16:07:51.14275122Z","level":"INFO","msg":"please access swagger docs","path":"http://localhost:8080/docs/index.html"}
...... output truncated ......
```

- How to build binary with `Makefile`
```
❯ make build                                                                                                                                                                                                                        
go build -o bin/manager cmd/run-services.go
```
Post, binary will be created in the `bin/` directory. 

- How to build docker image with `Makefile`
```
❯ make docker-build                                                                                                                                                                                                                 
docker build -t service:latest .
[+] Building 16.5s (15/15) FINISHED                                                                                                                                                                                
 => [internal] load build definition from Dockerfile                                                                                                                                                                    
 => => transferring dockerfile: 1.21kB                                                                                                                                                                                  
 => [internal] load .dockerignore                
 .
 .
 .
 
 ❯ make docker-push
```
