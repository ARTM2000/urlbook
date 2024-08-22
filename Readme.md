# URLBook
__URLBook__ is a sample project of URL-shortener platform with focus on backend side.

- [URLBook](#urlbook)
  - [System design](#system-design)
    - [Functional requirement](#functional-requirement)
    - [Non-Functional requirement](#non-functional-requirement)
    - [High-level Design Diagram](#high-level-design-diagram)
  - [Project architecture](#project-architecture)
  - [How to run the project](#how-to-run-the-project)
    - [Setup Running Environment](#setup-running-environment)
    - [Run the project](#run-the-project)
    - [Run the project for development](#run-the-project-for-development)
  - [Milestones](#milestones)
  - [Benchmark result](#benchmark-result)

## System design
### Functional requirement
- There should be an API (without OAuth action) that user can submit a valid URL for shortening.
- Shorted URL can have a system generated phrase or user customized value.
- The limitation for system generated phrase is up to 7 character
- The limitation for user customized phrase is up to 16 character
- The Shorted URL will not expire
- Shorted URL must redirect to the destination with status of 302
- With having shorted URL, can monitor the number of clicks and the device types which used for browse the link.

### Non-Functional requirement
- System should be highly available (99.9% uptime)
- 500k url submitted to the system per day
- 200M redirection request per month ( 200M / (30Day * 24h * 3600s) = 80tps )

### High-level Design Diagram
<div style='width: auto; max-width: 1000px; margin: 10px auto;'>
    <img src='docs/design.png' alt='high-level-system-design' />
</div>

## Project architecture
As considered, _Hexagonal Architecture_ selected for project structure to make the project flexible and isolated its parts.

Here are some links about the architecture:
 - https://medium.com/@pthtantai97/hexagonal-architecture-with-golang-part-1-7f82a364b29 
 - https://medium.com/@pthtantai97/hexagonal-architecture-with-golang-part-2-681ee2a0d780
 - https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/

## How to run the project
### Setup Running Environment
1. You at least need to have `docker` and `docker-compose` installed.
2. Run the following command in the project directory to setup the _environment variable_
   ```bash
   make prepare
   ```
   or if you don't have `make` installed, run the following
   ```bash
   bash ./scripts/prepare.bash
   ```
### Run the project
In order to run the project, use this command:
```bash
docker-compose -f ./deployments/docker-compose.yml --env-file ./.env up
```

### Run the project for development
For running in development mode, you should have  `golang >= 1.21` installed on your machine. After running the previous instruction
- Use the following to setup the require services
  ```bash
  docker-compose -f ./deployments/docker-compose.dev.yml --env-file ./.env up
  ```
- To run the project in watch mode, execute the following command: `make dev_server`


## Milestones
- [x] Create project system design
- [x] Setup required services with docker-compose
  - [x] Database (mysql)
  - [x] Cache (memcached)
- [x] Only submit a url and get a system generated short url
- [x] Redirect the system generated short-url to its original url with 302 HTTP status code
- [x] Submit a url with custom name for shortening
- [x] Bring caching mechanism
- [x] Add some tracking mechanism on urls
  - [x] Number of clicks with date
  - [x] The devices used to visit the link
  - [x] The IP address that clicks happened from 

## Benchmark result
[Apache benchmark (ab)](https://httpd.apache.org/docs/current/programs/ab.html) tool used for benchmarking the project, and how many request per second it can handle.

```txt
Server Software:        
Server Hostname:        localhost
Server Port:            3000

Document Path:          /L4ryUrj
Document Length:        0 bytes

Concurrency Level:      80
Time taken for tests:   2.510 seconds
Complete requests:      10000
Failed requests:        0
Non-2xx responses:      10000
Total transferred:      6450000 bytes
HTML transferred:       0 bytes
Requests per second:    3983.90 [#/sec] (mean)
Time per request:       20.081 [ms] (mean)
Time per request:       0.251 [ms] (mean, across all concurrent requests)
Transfer rate:          2509.39 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   1.1      0      26
Processing:     1   19  12.4     16     100
Waiting:        1   19  12.4     16     100
Total:          1   20  12.4     17     101

Percentage of the requests served within a certain time (ms)
  50%     17
  66%     21
  75%     23
  80%     26
  90%     35
  95%     47
  98%     59
  99%     67
 100%    101 (longest request)
```
