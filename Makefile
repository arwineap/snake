build:
	GOOS=js GOARCH=wasm go build -o web/snake.wasm main.go
	cp "${GOROOT}/misc/wasm/wasm_exec.js" ./wasm_exec.js

deploy:
	aws s3 sync web/ s3://snake-game-2021-8-17/ --acl public-read
