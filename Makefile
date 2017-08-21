run: data
	go run cmd/timer/main.go

data:
	go-bindata -nocompress -nometadata -pkg timer -o bindata.go assets/...
