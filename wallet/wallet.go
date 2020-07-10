package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checkSumlength = 4
	version        = byte(0x00) // hexadecimal representation of zero
)



type Wallet struct {
	//eliptic curve digital algorithm
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func ValidateAddres(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	checkSumFromAddress := pubKeyHash[len(pubKeyHash)-checkSumlength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checkSumlength]
	checkSum := CheckSum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(checkSum, checkSumFromAddress) == 0
}
func (w *Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionHash := append([]byte{version}, pubHash...)
	checksum := CheckSum(versionHash)

	fullHash := append(versionHash, checksum...)
	address := Base58Encode(fullHash)

	return address
}
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	return &Wallet{private, public}
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMd := hasher.Sum(nil)
	return publicRipMd
}

func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checkSumlength]
}
