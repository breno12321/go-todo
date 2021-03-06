# Golang-TODO App

## Development environment

âSpecs:

- RedisJSON - redislabs/rejson image
- Go: go1.16.6

## Step by step to execute ð¨

First of all, make sure your Postgres server is up and running - [Redis JSON](https://hub.docker.com/r/redislabs/rejson/) enviroments make its use easier. Then you can clone the repository with:  

```bash
git clone DESIRED_DIR_NAME git@github.com:breno12321/go-todo.git
```

Done clonning, you can now configure the env file, doing

```bash
vi .env.development
```

After setting up the enviroment you can run: 

```bash
go run cmd/go-todo/go-todo.go
```

for dev envs you can cd into `cmd/go-todo` and run `air` to watch file reloads - [Air lib](https://github.com/cosmtrek/air)

It will be running on `http://localhost:8080`.

## Folder Structure

```
âââ cmd
â   âââ go-todo
â       âââ go-todo.go
â       âââ tmp
â           âââ build-errors.log
â           âââ main
âââ go.mod
âââ go.sum
âââ internal
â   âââ items
â   â   âââ items.go
â   âââ structs
â       âââ structs.go
âââ readme.md
âââ src
â   âââ app
âââ tmp
    âââ build-errors.log
    âââ main
```
