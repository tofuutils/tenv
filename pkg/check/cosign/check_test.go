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

package cosigncheck_test

import (
	_ "embed"
	"testing"

	cosigncheck "github.com/tofuutils/tenv/v2/pkg/check/cosign"
	"github.com/tofuutils/tenv/v2/pkg/loghelper"
)

const (
	identity = "https://github.com/opentofu/opentofu/.github/workflows/release.yml@refs/heads/v1.6"
	issuer   = "https://token.actions.githubusercontent.com"
)

//go:embed testdata/tofu_1.6.0_SHA256SUMS
var data []byte

//go:embed testdata/tofu_1.6.0_SHA256SUMS.sig
var dataSig []byte

//go:embed testdata/tofu_1.6.0_SHA256SUMS.pem
var dataCert []byte

/*
 * no "t.Parallel()" on those tests (causes failures in cosign call)
 */

func TestCosignCheckCorrect(t *testing.T) { //nolint
	if err := cosigncheck.Check(data, dataSig, dataCert, identity, issuer, loghelper.InertDisplayer); err != nil {
		t.Error("Unexpected error :", err)
	}
}

func TestCosignCheckErrorCert(t *testing.T) { //nolint
	if cosigncheck.Check(data, dataSig, dataCert[1:], identity, issuer, loghelper.InertDisplayer) == nil {
		t.Error("Should fail on erroneous certificate")
	}
}

func TestCosignCheckErrorIdentity(t *testing.T) { //nolint
	if cosigncheck.Check(data, dataSig, dataCert, "me", issuer, loghelper.InertDisplayer) == nil {
		t.Error("Should fail on erroneous issuer")
	}
}

func TestCosignCheckErrorIssuer(t *testing.T) { //nolint
	if cosigncheck.Check(data, dataSig, dataCert, identity, "http://myself.com", loghelper.InertDisplayer) == nil {
		t.Error("Should fail on erroneous issuer")
	}
}

func TestCosignCheckErrorSig(t *testing.T) { //nolint
	if cosigncheck.Check(data, dataSig[1:], dataCert, identity, issuer, loghelper.InertDisplayer) == nil {
		t.Error("Should fail on erroneous signature")
	}
}
