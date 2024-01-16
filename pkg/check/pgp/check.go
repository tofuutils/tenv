package pgpcheck

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

func Check(data []byte, dataSig []byte, dataPublicKey []byte) error {
	pgpSignature := crypto.NewPGPSignature(dataSig)
	publicKeyObj, err := crypto.NewKeyFromArmored(string(dataPublicKey))
	if err != nil {
		return err
	}

	signingKeyRing, err := crypto.NewKeyRing(publicKeyObj)
	if err != nil {
		return err
	}

	message := crypto.NewPlainMessage(data)
	return signingKeyRing.VerifyDetached(message, pgpSignature, crypto.GetUnixTime())
}
