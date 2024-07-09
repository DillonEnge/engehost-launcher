run:
	@go run cmd/main.go

build:
	@go build -o dist/darwin/Engehost\ Launcher.app/Contents/MacOS/engehost_launcher cmd/main.go

build_win:
	@go-winres simply --icon assets/icon.png
	@GOOS=windows GOARCH=amd64 go build -o dist/windows/EngehostLauncher.exe cmd/main.go

open:
	@open dist/darwin/Engehost\ Launcher.app
