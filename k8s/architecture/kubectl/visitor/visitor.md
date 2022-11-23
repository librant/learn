
```shell
$ go run .
2022/11/15 08:56:09 visitor.go:37: In Visitor1 before fn
2022/11/15 08:56:09 visitor.go:19: In VisitorList before fn
2022/11/15 08:56:09 visitor.go:51: In Visitor2 before fn
2022/11/15 08:56:09 visitor.go:70: In Visitor3 before fn
2022/11/15 08:56:09 main.go:22: In visitFunc
2022/11/15 08:56:09 visitor.go:74: In Visitor3 after fn
2022/11/15 08:56:09 visitor.go:55: In Visitor2 after fn
2022/11/15 08:56:09 visitor.go:23: In VisitorList after fn
2022/11/15 08:56:09 visitor.go:41: In Visitor1 after fn
```