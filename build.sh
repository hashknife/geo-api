#!/bin/sh

make build
docker build -t geo-api:latest .
