/*
 *
 * Copyright 2024 tofuutils authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package pgpcheck

import (
	"errors"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

var (
	ErrCheck    = errors.New("invalid pgp signature")
	ErrNoKey    = errors.New("no valid pgp public keys found")
	ErrKeyBlock = errors.New("invalid pgp key block")
)

func Check(data []byte, dataSig []byte, dataPublicKey []byte) error {
	keys, err := extractKeys(dataPublicKey)
	if err != nil {
		return err
	}

	pgpSignature := crypto.NewPGPSignature(dataSig)
	message := crypto.NewPlainMessage(data)

	for _, key := range keys {
		keyRing, err := crypto.NewKeyRing(key)
		if err != nil {
			continue
		}

		if err = keyRing.VerifyDetached(message, pgpSignature, crypto.GetUnixTime()); err == nil {
			return nil
		}
	}

	return ErrCheck
}

func extractKeys(dataPublicKey []byte) ([]*crypto.Key, error) {
	keyStr := string(dataPublicKey)
	keySep := "-----BEGIN PGP PUBLIC KEY BLOCK-----"

	capacity := strings.Count(keyStr, keySep)

	keys := make([]*crypto.Key, 0, capacity)

	for part := range strings.SplitSeq(keyStr, keySep) {
		if strings.TrimSpace(part) == "" {
			continue
		}

		armoredKey := keySep + part
		publicKeyObj, err := crypto.NewKeyFromArmored(armoredKey)
		if err != nil {
			return nil, ErrKeyBlock
		}

		keys = append(keys, publicKeyObj)
	}

	if len(keys) == 0 {
		return nil, ErrNoKey
	}

	return keys, nil
}
