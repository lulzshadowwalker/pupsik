CMD := $(firstword $(MAKECMDGOALS))
ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
%::
	@true

.PHONY: dev run install build goose-up goose-up-by-one goose-up-to goose-down goose-down-to goose-redo goose-reset goose-status goose-version goose-create goose-fix goose-validate

dev:
	@air

run: build
	@./bin/pupsik

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install --global tailwindcss
	@npm install --save-dev daisyui@latest

build:
	@tailwindcss -i view/css/app.css -o public/styles.css
	@templ generate view
	@go build -o bin/pupsik main.go

GOOSE_DIR="./database/migrations/"
# goose up
goose-up:
	@goose --dir ${GOOSE_DIR} up

# goose up-by-one
goose-up-by-one:
	@goose --dir ${GOOSE_DIR} up-by-one

# goose up-to VERSION
goose-up-to:
	@goose --dir ${GOOSE_DIR} up-to $(ARGS)

# goose down
goose-down:
	@goose --dir ${GOOSE_DIR} down

# goose down-to VERSION
goose-down-to:
	@goose --dir ${GOOSE_DIR} down-to $(ARGS)

# goose redo
goose-redo:
	@goose --dir ${GOOSE_DIR} redo

# goose reset
goose-reset:
	@goose --dir ${GOOSE_DIR} reset

# goose status
goose-status:
	@goose --dir ${GOOSE_DIR} status

# goose version
goose-version:
	@goose --dir ${GOOSE_DIR} version

# goose create NAME [sql|go] 
goose-create:
	@goose --dir ${GOOSE_DIR} -s create $(ARGS)

# goose fix
goose-fix:
	@goose --dir ${GOOSE_DIR} fix

# goose validate
goose-validate:
	@goose --dir ${GOOSE_DIR} validate

# Run Jet migrations
jet-migrations:
	@go run cmd/migrations/main.go
