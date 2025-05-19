# How to use protobuf
1. install
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

2. generate from ./proto dir (the *.proto has saved in the dir)
```
cd proto
protoc --go_out=../ *.proto
```