# tx-track-engine-go

Track and post-process transactions from different blockchains.

## Languages and Frameworks

- [Go](https://go.dev) - Build fast, reliable, and efficient software at scale

## Local Environment Setup

### Golang Environment Variable Setup (ZSH)

```bash
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOPATH
export PATH=$PATH:$GOROOT/bin
export PATH=$PATH:$GOPATH/bin
```

### Postgresql Docker Setup

```bash
# Setup container running postgres database
docker pull postgres:14.1

docker run -d --name postgres -e POSTGRES_USERNAME=root -e POSTGRES_PASSWORD=root -v ${HOME}/<path>/:/var/lib/postgresql/data -p 5432:5432 postgres:14.1

docker exec -it postgres bash

# Setup superuser root with password root
root@e0406f495e62:/ su - postgres

postgres@e0406f495e62:~$ createuser --interactive --pwprompt
```

### Migration Setup

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

vim .zshrc
export PATH=$PATH:$GOPATH/bin
```
