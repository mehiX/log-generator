# LOGS generator

Generate a lot of message logs in a given period. This can be used for load testing.

## Requirements

- Go >=1.19
- Docker

If running in Docker then you don't need to install Go

## Run

See the help for parameter details:

```
go run ./... -h
```

Example of run: 

```
# Run for 1 second with 13 concurrent generators and print only stats
go run ./... -run 1000ms -g 13 -stats

# Run for 1 second with 10 concurrent generators and generate actual messages
go run ./... -run 1000ms -g 10
```

## Run with Docker

```
# Run for 1 second with 13 concurrent generators and print only stats
docker build -t test . && docker run -ti --rm test -run 1000ms -g 13 -stats
```

For now the input file is copied at build time. The program can be easily changed to mount a volume at container runtime.