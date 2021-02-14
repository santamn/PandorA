WinCC := x86_64-w64-mingw32-gcc

macos: # Not for cross-compile
	cd ./cmd/pandora && GOOS=darwin GOARCH=amd64 go build
	mv ./cmd/pandora/pandora ./bin/macos/
	cd ./cmd/form && GOOS=darwin GOARCH=amd64 go build
	mv ./cmd/form/form ./bin/macos/

linux: # Not for cross-compile
	cd ./cmd/pandora && GOOS=linux GOARCH=amd64 go build
	mv ./cmd/pandora/pandora ./bin/linux/
	cd ./cmd/form && GOOS=linux GOARCH=amd64 go build
	mv ./cmd/form/form ./bin/linux/

windows: # For cross-compile (Mac -> Windows)
	cd ./cmd/pandora && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=$(WinGCC) go build -ldflags "-H=windowsgui"
	mv ./cmd/pandora/pandora.exe ./bin/windows/
	cd ./cmd/form && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=$(WinGCC) go build -ldflags "-H=windowsgui"
	mv ./cmd/form/form.exe ./bin/windows/