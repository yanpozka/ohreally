## Oh Really!
A tool to identify this kind of code:
```go
if err != nil {
    return err
}
return nil
```

### Install:

```
go get github.com/yanpozka/ohreally
ohreally -file example/example.go
```
