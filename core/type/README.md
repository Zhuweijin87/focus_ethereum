## 核心交易模块 

### transaction.go 
这里主要包含了交易结构数据定义 

```go
type Transaction struct {
	data txdata   // 数据
	
	hash atomic.Value  
	size atomic.Value
	from atomic.Value
}
```

具体数据：
```go
type txdata struct {
	AccountNonce uint64 
	Price        *big.Int  // 支付交易的价格
	GasLimit     *big.Int    // gas 上限
	Recipient    *common.Address  // 接收地址，如果是合约交易为nil
	Amount       *big.Int  // 金额    
	Payload      []byte  // 数据  

	// 签名使用的数据 
	V *big.Int 
	R *big.Int 
	S *big.Int

	Hash *common.Hash // 转化为Json数据时，需要hash
}
``` 

主要函数  
+ NewTransaction(nonce uint64, to common.Address, amount, gasLimit, gasPrice *big.Int, data []byte) *Transaction