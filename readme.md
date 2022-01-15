# Golang-TODO App

## Development environment

⚙Specs:

- RedisJSON - redislabs/rejson image
- Go: go1.16.6

## Step by step to execute 💨

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
├── cmd
│   └── go-todo
│       ├── go-todo.go
│       └── tmp
│           ├── build-errors.log
│           └── main
├── go.mod
├── go.sum
├── internal
│   ├── items
│   │   └── items.go
│   └── structs
│       └── structs.go
├── readme.md
├── src
│   └── app
└── tmp
    ├── build-errors.log
    └── main
```
