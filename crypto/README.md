## 加密模块 

以太坊的安全都是基于加密算法的，该模块封装的大多使用的加密算法:  

### crypto.go

+ GenerateKey() (*ecdsa.PrivateKey, error) : 生产私钥  
通过私钥获取公钥： 
```go
    privkey, err := GenerateKey()
    pubkey := privkey.PublicKey
```

+ ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) 校验签名  

+ PubkeyToAddress(p ecdsa.PublicKey) common.Address 将公钥转化地址(取公钥后20位) 

+ Keccak256(data ...[]byte) []byte  返回一个32位hash  

+ Keccak512(data ...[]byte) []byte  返回一个64位hash  

+ Keccak256Hash(data ...[]byte) (h common.Hash)  

