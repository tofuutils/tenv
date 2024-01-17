/*
 *
 * Copyright 2024 gotofuenv authors.
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

package cosigncheck_test

import (
	_ "embed"
	"testing"

	cosigncheck "github.com/dvaumoron/gotofuenv/pkg/check/cosign"
)

//go:embed tofu_1.6.0_arm64.deb
var data []byte

//go:embed tofu_1.6.0_arm64.deb.sig
var dataSig []byte

//go:embed tofu_1.6.0_arm64.deb.pem
var dataKey []byte

func TestCosignCheckCorrect(t *testing.T) {
	if err := cosigncheck.Check(data, dataSig, dataKey); err != nil {
		t.Error("Unexpected error : ", err)
	}
}

func TestCosignCheckErrorKey(t *testing.T) {
	if cosigncheck.Check(data, dataSig, dataKey[1:]) == nil {
		t.Error("Should fail on erroneus public key")
	}
}

func TestCosignCheckErrorSig(t *testing.T) {
	if cosigncheck.Check(data, dataSig[1:], dataKey) == nil {
		t.Error("Should fail on erroneus signature")
	}
}
