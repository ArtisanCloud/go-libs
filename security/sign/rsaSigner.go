package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"github.com/ArtisanCloud/PowerLibs/v3/security/encryption"
)

type RSASigner struct {
	RSAEncryptor *encryption.RSAEncryptor
}

func NewRSASigner(hash crypto.Hash) (signer *RSASigner, err error) {

	encryptor, err := encryption.NewRSAEncryptor(hash)
	if err != nil {
		return nil, err
	}
	signer = &RSASigner{
		RSAEncryptor: encryptor,
	}

	return signer, err
}

// ----------------------------------------------------------------------------------------------------
// Before signing, we need to hash our sms
// The hash is what we actually sign
// ----------------------------------------------------------------------------------------------------

// Sign
// use hash method to make an actual sign
func (signer *RSASigner) Sign(msg []byte) (msgHashSum []byte, err error) {

	msgHash := signer.RSAEncryptor.Hash.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		return nil, err
	}
	msgHashSum = msgHash.Sum(nil)

	return msgHashSum, err
}

// ----------------------------------------------------------------------------------------------------
// In order to generate the signature, we provide a random number generator,
// our private key, the hashing algorithm that we used, and the hash sum
// of our sms
// ----------------------------------------------------------------------------------------------------

// GenerateSignaturePSS
// use this signed sms to create a digital signature via private key
// digest has been sum hashed
func (signer *RSASigner) GenerateSignaturePSS(digest []byte) (signature []byte, err error) {
	hash := signer.RSAEncryptor.Hash
	signature, err = rsa.SignPSS(rand.Reader, signer.RSAEncryptor.PrivateKey, hash, digest, nil)
	return signature, err
}

// GenerateSignaturePKCS1v15
// use this signed sms to create a digital signature via private key
// digest has been sum hashed
func (signer *RSASigner) GenerateSignaturePKCS1v15(digest []byte) (signature []byte, err error) {
	hash := signer.RSAEncryptor.Hash
	signature, err = rsa.SignPKCS1v15(rand.Reader, signer.RSAEncryptor.PrivateKey, hash, digest)
	return signature, err
}

// ----------------------------------------------------------------------------------------------------
// To verify the signature, we provide the public key, the hashing algorithm
// the hash sum of our sms and the signature we generated previously
// there is an optional "options" parameter which can omit for now
// ----------------------------------------------------------------------------------------------------

func (signer *RSASigner) VerifySignPSS(digest []byte, signature []byte) (err error) {
	err = rsa.VerifyPSS(signer.RSAEncryptor.PublicKey, signer.RSAEncryptor.Hash, digest, signature, nil)
	return err
}

func (signer *RSASigner) VerifySignPKCS1v15(digest []byte, signature []byte) (err error) {

	err = rsa.VerifyPKCS1v15(signer.RSAEncryptor.PublicKey, signer.RSAEncryptor.Hash, digest, signature)
	return err
}
