## 钱包账号 

### account.go
```go
type Account struct {
	Address common.Address  // 账号地址 
	URL     URL             // 
}
```

```go
type URL struct {
	Scheme string // 
	Path   string // 
}
```
URL 是对节点的解析  
例如: enode://a979fb575495b8d6db44f750317d0f4622bf4c2aa3365d6af7c284339968eef29b69ad0dce72a4d8db5ebb4968de0e3bec910127f134779fbcb0cb6d3331163c@52.16.188.185:30303  
转化为URL后，如下：  
Scheme:enode 
Path:a979fb575495b8d6db44f750317d0f4622bf4c2aa3365d6af7c284339968eef29b69ad0dce72a4d8db5ebb4968de0e3bec910127f134779fbcb0cb6d3331163c@52.16.188.185:30303 