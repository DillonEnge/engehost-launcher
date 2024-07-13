run:
	@go run cmd/main.go

build_darwin:
	@go build -o dist/darwin/Engehost\ Launcher.app/Contents/MacOS/engehost_launcher cmd/main.go

build_win:
	@go-winres simply --icon assets/icon.png --file-version git-tag --admin
	@mv rsrc_windows_* cmd/
	@GOOS=windows GOARCH=amd64 go build -o dist/windows/EngehostLauncher.exe cmd/main.go
	@rm cmd/rsrc_windows_*

open:
	@open dist/darwin/Engehost\ Launcher.app
