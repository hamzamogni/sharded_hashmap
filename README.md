# Sharded Hashmap Implementation 

## Introduction

This is an implementation of a sharded hashmap, sharded over network of multiple nodes.

## Design

### Master Node

The master is responsible for the following tasks:
- Registering nodes
- Distributing keys to nodes
- Receiving requests from clients and forwarding them to the appropriate node
- Receiving responses from nodes and forwarding them to the appropriate client

The master node decides which shard node to store the key in by calculating a hashed value of the key and then taking the modulo of the number of nodes,
this means that number of nodes should be fixed and known in advance. Dynamic addition of nodes is not supported (yet).

### Shard Node

The shard node is responsible for the following tasks:
- Storing a subset of the keys
- Receiving requests from the master node and responding to them


Data is kept in an in-memory hashmap, and is not persisted to disk.

## Usage

A docker-compose file is provided to run the master and shard nodes. To run the system, simply run `docker-compose up` in the root directory of the project.

The master node listens on port 8080, and the shard nodes listen on ports 9091, 9092, 9093, etc.

The default docker-compose file runs 4 shards, but this can be changed by modifying the `docker-compose.yml` file.

### CLI Client

A CLI is provided to interact with the master node. 

```golang
go run main.go --help
```
