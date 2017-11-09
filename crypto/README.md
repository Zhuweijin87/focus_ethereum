## 加密模块 

以太坊的安全都是基于加密算法的，该模块封装的大多使用的加密算法:  

### crypto.go

+ GenerateKey() (*ecdsa.PrivateKey, error) : 生产私钥  
通过私钥获取公钥： 
```go
    privkey, err := crypto.GenerateKey()
    pubkey := privkey.PublicKey
    // 或者 pubKey := crypto.
```

+ ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) 校验签名  

+ PubkeyToAddress(p ecdsa.PublicKey) common.Address 将公钥转化地址(取公钥后20位) 

+ Keccak256(data ...[]byte) []byte  返回一个32位hash  

+ Keccak512(data ...[]byte) []byte  返回一个64位hash  

+ Keccak256Hash(data ...[]byte) (h common.Hash)  

+ CreateAddress(b common.Address, nonce uint64): 根据给定的地址和难度系数算出以太坊地址 

+ Sign(hash []byte, priv *ecdsa.PrivateKey) ([]byte, error): 签名 
例子:  
```go
    hash := "12345abcdef2345612345abcdef23456"
    sign, err := crypto.Sign(hash, privKey)
```

+ SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error): 签名之后的数据转化为公钥 

关于椭圆曲线加密说明： http://8btc.com/article-138-1.html  
+ secp256k1: 椭圆曲线密码 
+ Keccak256: 基于SHA-256的散列算法  
+ Keccak512: 基于SHA-512的散列算法  