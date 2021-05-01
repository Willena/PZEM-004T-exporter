# PZEM-004T-exporter

A prometheus exporter for PZEM-004T power meter written in Golang

## Build instructions

1. Install Golang
2. clone the repository
3. `go build -o PZEM_004T_exporter .`

## Usage

```
Usage of PZEM_004T_exporter:
  -host string
    	Hostname to bind to (default "0.0.0.0")
  -port int
    	Port to listen request on (default 2112)
  -resetEnergy
    	Should the energy value be reset at start
  -serialPort string
    	Serial port used to communicate with PZEM-004T
```

## Docker

This program is also available as a Docker container. The container should be privileged to allow acces to the serial 
port. Another solution would be to write a udev rule to allow to read and write to them.

```
docker run --rm --privileged -v /dev/ttyS1:/dev/ttyS1 -p 2112:2112 gillena/pzem-exporter -serialPort /dev/ttyS1
```