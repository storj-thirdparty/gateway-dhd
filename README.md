  # gateway-dhd

[![Go Report Card](https://goreportcard.com/badge/github.com/storj-thirdparty/gateway-dhd)](https://goreportcard.com/report/github.com/storj-thirdparty/gateway-dhd)


## Overview

Dial Home Device. Small gateway for protocol translation.

It is an HTTP REST Server to stream the response from the decentralized Storj network to the client. Written in Golang.

```
Usage:
  gateway-dhd [command] <flags>

Available Commands:
  help        Help about any command
  start       start the REST server

Flags:
  -h, --help          help for gateway-dhd
  -p, --port string   Port number of the REST server (default "8080")
```  
  
```start``` - Start the REST server at the specified port number specified in the ```port``` flag (default port: ```8080```).


## Requirements and Install
To build from scratch, [install the latest Go](https://golang.org/doc/install#install).

> Note: Ensure go modules are enabled (GO111MODULE=on)

#### Option #1: clone this repo (most common)
To clone the repo
```
git clone https://github.com/storj-thirdparty/gateway-dhd.git
```
Then, build the project using the following:
```
cd gateway-dhd
go build
```
#### Option #2: go get into your gopath
To download the project inside your GOPATH use the following command:
```
go get github.com/storj-thirdparty/gateway-dhd
```
## Run
Once you have built the project run the following commands as per your requirement:

##### Start the REST server at default port 8080
```
$ ./gateway-dhd start
```
##### Start the REST server at custom port
```
$ ./gateway-dhd start --port <custom_port_number>
```
##### Get help
```
$ ./gateway-dhd --help
```

## Testing
* The project has been tested on the following operating systems:
```
	* Windows
		* Version: 10 Pro
		* Processor: Intel(R) Core(TM) i3-5005U CPU @ 2.00GHz 2.00GHz

	* macOS Catalina
		* Version: 10.15.4
		* Processor: 2.5 GHz Dual-Core Intel Core i5

	* ubuntu
		* Version: 16.04 LTS
		* Processor: AMD A6-7310 APU with AMD Radeon R4 Graphics × 4
```	