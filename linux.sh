#!/bin/bash
GOOS=linux GOARCH=amd64 go build && tar -zcvf jetson-monitor.tar.gz jetson-monitor