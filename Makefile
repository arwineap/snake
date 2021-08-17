build:
	GOOS=js GOARCH=wasm go build -o snake.wasm main.go
