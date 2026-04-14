swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/gin-investment-tracker/main.go -o docs

test:
	go test -count=1 ./...

install-hooks:
	cp scripts/hooks/pre-push .git/hooks/pre-push
	chmod +x .git/hooks/pre-push
	@echo "Git hooks installed. Tests will run before every push to main."