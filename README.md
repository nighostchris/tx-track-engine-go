# tx-track-engine-go

Track and post-process transactions from different blockchains.

## Languages and Frameworks

- [Go](https://go.dev) - Build fast, reliable, and efficient software at scale

## How To Run the Program

```bash
go run main.go
```

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

## Program Output Demonstration

```bash
[Database Migration] Done
[EVM Node RPC] GetBlockByNumber - 9160005
[EVM Node RPC] GetBlockByNumber - 9160003
[EVM Node RPC] GetBlockByNumber - 9160004
[EVM Node RPC] GetBlockByNumber - 9160002
[EVM Node RPC] GetBlockByNumber - 9160001
[EVM Node RPC] GetBlockByNumber -  Block #9160004 have 0 transactions
[EVM Process Block] Finished processing block 9160004 and found 0 interested transaction(s)
[EVM Node RPC] GetBlockByNumber -  Block #9160003 have 0 transactions
[EVM Node RPC] GetBlockByNumber -  Block #9160002 have 0 transactions
[EVM Process Block] Finished processing block 9160003 and found 0 interested transaction(s)
[EVM Process Block] Finished processing block 9160002 and found 0 interested transaction(s)
[EVM Node RPC] GetBlockByNumber -  Block #9160005 have 1 transactions
[EVM Process Block] Finished processing block 9160005 and found 0 interested transaction(s)
[EVM Node RPC] GetBlockByNumber -  Block #9160001 have 5 transactions
[EVM Process Block] Finished processing block 9160001 and found 0 interested transaction(s)
[EVM Node RPC] GetBlockByNumber - 9160010
[EVM Node RPC] GetBlockByNumber - 9160007
[EVM Node RPC] GetBlockByNumber - 9160006
[EVM Node RPC] GetBlockByNumber - 9160008
[EVM Node RPC] GetBlockByNumber - 9160009
[EVM Node RPC] GetBlockByNumber -  Block #9160007 have 0 transactions
[EVM Process Block] Finished processing block 9160007 and found 0 interested transaction(s)
[EVM Node RPC] GetBlockByNumber -  Block #9160010 have 2 transactions
[EVM Process Block] Finished processing block 9160010 and found 0 interested transaction(s)
[EVM Node RPC] GetBlockByNumber -  Block #9160008 have 2 transactions
[EVM Process Block] Finished processing block 9160008 and found 0 interested transaction(s)
[EVM Node RPC] GetBlockByNumber -  Block #9160009 have 0 transactions
[EVM Process Block] Finished processing block 9160009 and found 0 interested transaction(s)
[EVM Node RPC] GetBlockByNumber -  Block #9160006 have 3 transactions
[EVM Process Block] Finished processing block 9160006 and found 0 interested transaction(s)
```
