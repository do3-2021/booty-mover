# Booty-mover

A server management bot written in go.

## Environment variables

- BOT_TOKEN : token of the discord bot.
- POSTGRES : uri to connect to the postgres database eg `postgres://postgres:postgres@localhost/postgres?sslmode=disable`.

## How to launch

- Start a PostgerSQL server
- Set the environment variables
- Execute the default executable (`make && ./booty-mover` or `go run .`)

## Dependencies

- make
- go
- pandoc for manual generation

## make commands

- `make all` : builds for Windows, linux generic and Ubuntu/debian (deb), builds the manuals and put everything in the `publish` folder
- `make` : builds for your current platform
