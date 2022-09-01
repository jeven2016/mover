### init module
```shell
go mod init my
```

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

### list all mod
```shell
 go list -m all
```

### check versions 
```shell
go list -m -version pkgName
```
* go mod init 创建一个新模块，初始化描述它的go.mod文件。
* go buil，go test和其他程序包构建命令根据需要向go.mod添加新的依赖项。
* go list -m all打印当前模块的依赖关系。
* go get更改所需依赖的版本（或添加新的依赖）。
* go mod tidy删除未使用的依赖项。
* go clean -modcache  清除所有mod缓存
* go clean -i github.com/mix-go/dotenv... 执行删除编译后的package目录
  * 请一定要包含三个点号 ... ，这样就不会递归删除子package，如本例中的 gore/gocode。-i 参数表示删除由 go install 所创建的archive或二进制文档。






