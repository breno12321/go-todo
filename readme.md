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
go run main.go
```

It will be running on `http://localhost:8080`.

## Folder Structure

```
├── components
│   └── items
│       ├── api.go
│       ├── model.go
│       └── routes.go
├── go.mod
├── go.sum
├── main.go
└── readme.md
```
