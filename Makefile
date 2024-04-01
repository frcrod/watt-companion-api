generate:
	@go mod tidy
	@templ generate
	@tailwind -c .\config\tailwind.config.js -o ./web/static/css/output.css --minify
	@go run cmd/main.go

run: 
	@go run cmd/main.go
