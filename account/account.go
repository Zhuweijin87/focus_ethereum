package main

import (
	"fmt"
	"io/ioutil"
	"os"
	_ "strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	_ "github.com/ethereum/go-ethereum/common"
	_ "github.com/ethereum/go-ethereum/event"
)

func tmpKeyStore(encrypt bool) (string, *keystore.KeyStore) {
	d, err := ioutil.TempDir("", "eth-keystore-test") // 创建临时文件
	if err != nil {
		fmt.Println("fail to TempDir: ", err)
		return "nil", nil
	}

	new := keystore.NewPlaintextKeyStore
	if encrypt {
		new = func(kd string) *keystore.KeyStore {
			return keystore.NewKeyStore(kd, 2, 1) // 路径， scryptN, scryptP 生产KeyStore
		}
	}
	return d, new(d)
}

// 使用账号签名
func SignHash(acct accounts.Account, ks *keystore.KeyStore) {
	err := ks.Unlock(acct, "foo")
	if err != nil {
		fmt.Println(err)
		return
	}

	testSigData := make([]byte, 32)

	_, err = ks.SignHash(accounts.Account{Address: acct.Address}, testSigData)
	if err != nil {
		fmt.Println("Sign fail")
		return
	}

	fmt.Printf("sign: %x\n", testSigData)
}

// 用账号密码签名
func SignHashWithPass(acc accounts.Account, pass string, ks *keystore.KeyStore) {
	err := ks.Unlock(acc, "foo")
	if err != nil {
		fmt.Println(err)
		return
	}

	testSigData := make([]byte, 32)
	_, err = ks.SignHashWithPassphrase(acc, pass, testSigData)
	if err != nil {
		fmt.Println("SignHashWith Pass :", err)
		return
	}

	fmt.Printf("SIgn pass hash: %x\n", testSigData)
}

func main() {
	dir, ks := tmpKeyStore(true)
	defer os.RemoveAll(dir)

	fmt.Println(dir)

	a, err := ks.NewAccount("foo") // 创建账号 foo - 密码
	if err != nil {
		fmt.Println("NewAccount: ", err)
		return
	}
	fmt.Printf("Type a: %T\n", a)
	fmt.Println("account url path: ", a.URL.Path)
	fmt.Printf("account address: %x\n", a.Address)

	SignHash(a, ks)
	SignHashWithPass(a, "foo", ks)
}
