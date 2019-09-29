run:
	go run cmd/commandline/main.go ${ARGS}

aws-lambda-build:
	GOOS=linux go build -o main cmd/aws-lambda/main.go
	zip -j aws-lambda-weatherclockapi.zip main
	rm main

restapi-server-build:
	GOOS=linux go build -o weatherclock-api cmd/rest-api/main.go

cmdtool-build:
	go build -o weatherclock-cli cmd/commandline/main.go
