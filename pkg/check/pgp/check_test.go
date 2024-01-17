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

package pgpcheck_test

import (
	_ "embed"
	"testing"

	pgpcheck "github.com/dvaumoron/gotofuenv/pkg/check/pgp"
)

//go:embed terraform_1.6.6_SHA256SUMS
var data []byte

//go:embed terraform_1.6.6_SHA256SUMS.sig
var dataSig []byte

//go:embed hashicorp-pgp-key.txt
var dataKey []byte

func TestPgpCheckCorrect(t *testing.T) {
	if err := pgpcheck.Check(data, dataSig, dataKey); err != nil {
		t.Error("Unexpected error : ", err)
	}
}

func TestPgpCheckErrorKey(t *testing.T) {
	if pgpcheck.Check(data, dataSig, dataKey[1:]) == nil {
		t.Error("Should fail on erroneus public key")
	}
}

func TestPgpCheckErrorSig(t *testing.T) {
	if pgpcheck.Check(data, dataSig[1:], dataKey) == nil {
		t.Error("Should fail on erroneus signature")
	}
}
