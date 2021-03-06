package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testKey       = "30770201010420486227dccfc31c508b37ca4f3fb3fcbb851b5512159aff44bf913c2ad7db4381a00a06082a8648ce3d030107a14403420004ab226dd60c8dc077ccf4d31ef5b7935ab28c9a2222bef44adba15ec59ab16fa23ff7c820486a95bdbdde8c6b3e5cc0ad7ee5c2f83621d91cc27030fa8687c495"
	testPayload   = "00e84fadf17df0e08115e5694ed63a3e26e699c3d4239ebe9f66e925d5dda319"
	testSignature = "1e40ddf7f0bcc1a102a772b92cedaa548fa8ef7ec4da7bfe627b2e832e4793a1c57f1e3e7c54a6b0e87bf0902eb3caa748e96f7deacb01e87cad39287203c39c"
)

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

func TestSign(t *testing.T) {
	s := Sign(testPayload, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign should return a hex encoded string, got %s", s)
	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("blah")
	if err == nil {
		t.Error("restoreBigInts should return error when payload is not hex")
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{testPayload, true},
		{"00e84fadf17df0e08115e5694ed63a3e26e699c3d4239ebe9f66e925d5dda310", false},
	}

	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSignature, w.Address, tc.input)
		if ok != tc.ok {
			t.Error("Verify could not verify testSignature and testPayload")
		}
	}
}

func TestWallet(t *testing.T) {
	t.Run("New wallet is created", func(t *testing.T) {
		f = fakeLayer{
			fakeHasWalletFile: func() bool {
				return false
			},
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New wallet should return a new wallet instance")
		}
	})

	t.Run("Wallet is restored", func(t *testing.T) {
		f = fakeLayer{
			fakeHasWalletFile: func() bool {
				return true
			},
		}
		w = nil
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New wallet should return a new wallet instance")
		}
	})
}

func TestHasWalletFile(t *testing.T) {
	f := layer{}
	b := f.hasWalletFile()
	if b == true {
		t.Error("HasWalletFile should return false ")
	}
}

func TestWriteFile(t *testing.T) {
	f := layer{}
	err := f.writeFile("", []byte(""), 0644)
	if err == nil {
		t.Error("WriteFile should return error")
	}
}

func TestReadFile(t *testing.T) {
	f := layer{}
	_, err := f.readFile("test")
	if err == nil {
		t.Error("ReadFile should return error")
	}
}
