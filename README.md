### add dependency
```
 go get gopkg.in/ini.v1
```

### run with arguments
```shell
 ./main -from=hello -to=where
```

### download 
```shell
go mod download
```

### build 
```shell
# Mac
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build filename.go
 
# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build filename.go
```