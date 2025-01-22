# General
### Compile proto
```bash
protoc --go_out=. --go-grpc_out=. ./scooter.proto 
```

# Scooter Server
### Build docker image
```bash
docker build . -t scooter-server:0.2
```
### Run server container
```bash
docker run -p 50051:50051 --name scooter-server scooter-server:0.2
```

### Docker compose
Change to `/src/docker/` directory and use the following commands:
```bash
docker-compose up -d # to start all containers and load balancer
docker-compose down # to delete alll containers
```