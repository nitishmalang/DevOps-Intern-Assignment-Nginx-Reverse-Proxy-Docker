#!/bin/bash
echo "Testing Service 1 (Golang)..."
curl -s http://localhost:8080/service1/ping > /dev/null && echo "Service 1 /ping OK" || echo "Service 1 /ping failed"
curl -s http://localhost:8080/service1/hello > /dev/null && echo "Service 1 /hello OK" || echo "Service 1 /hello failed"
echo "Testing Service 2 (Python)..."
curl -s http://localhost:8080/service2/ping > /dev/null && echo "Service 2 /ping OK" || echo "Service 2 /ping failed"
curl -s http://localhost:8080/service2/hello > /dev/null && echo "Service 2 /hello OK" || echo "Service 2 /hello failed"
echo "Checking Nginx logs..."
docker-compose logs nginx | tail -n 5