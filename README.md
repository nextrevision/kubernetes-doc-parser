# kubernetes-doc-parser

This tool really has no use beyond being a quick script for [kubernetes-docset-generator](https://github.com/nextrevision/kubernetes-docset-generator)

## Installing

```
go get -u github.com/nextrevision/kubernetes-doc-parser
```

## Building

```
go get -u github.com/kardianos/govendor
govendor sync
GOOS=darwin GOARCH=amd64 go build
```
