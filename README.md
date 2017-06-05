## Oh Really!
A tool to identify this kind of code:
```go
if err != nil {
    return err
}
return nil
```
that can be written as just:
```go
return err
```

### Install and try it:

```bash
go get -v -u github.com/yanpozka/ohreally

cd $GOPATH github.com/yanpozka/ohreally/example

ohreally -file example.go
```
