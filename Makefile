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
