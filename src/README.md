# General
### Compile proto
```bash
protoc --go_out=. --go-grpc_out=. ./scooter.proto 
protoc --go_out=. --go-grpc_out=. ./multipaxos.proto 
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


### etcd
Read all keys:
```bash
export ETCDCTL_API=3
etcdctl get "" --prefix --keys-only
```

Read all registered servers:
```bash
etcdctl get "/servers/" --prefix
```