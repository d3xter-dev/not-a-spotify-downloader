build:
	GOOS=windows GOARCH=amd64 go build -o not-a-spotify-downloader.exe ./cmd/main.go
	GOOS=linux GOARCH=amd64 go build  -o not-a-spotify-downloader-linux-amd64 ./cmd/main.go
	go build -o not-a-spotify-downloader-linux-arm ./cmd/main.go
