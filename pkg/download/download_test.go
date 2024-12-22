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

package download_test

import (
	"testing"

	"github.com/tofuutils/tenv/v4/pkg/download"
)

func TestURLTransformer(t *testing.T) {
	t.Parallel()

	urlTransformer := download.NewURLTransformer("https://releases.hashicorp.com", "http://localhost:8080")

	value, err := urlTransformer("https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_linux_386.zip")
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	if value != "http://localhost:8080/terraform/1.7.0/terraform_1.7.0_linux_386.zip" {
		t.Error("Unexpected result, get :", value)
	}
}

func TestURLTransformerPrefix(t *testing.T) {
	t.Parallel()

	urlTransformer := download.NewURLTransformer("https://github.com", "https://go.dev")

	initialValue := "https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_darwin_amd64.zip"
	value, err := urlTransformer(initialValue)
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	if value != initialValue {
		t.Error("Unexpected result, get :", value)
	}
}

func TestURLTransformerSlash(t *testing.T) {
	t.Parallel()

	urlTransformer := download.NewURLTransformer("https://releases.hashicorp.com/", "http://localhost")

	value, err := urlTransformer("https://releases.hashicorp.com/terraform/1.7.0/terraform_1.7.0_darwin_amd64.zip")
	if err != nil {
		t.Fatal("Unexpected error :", err)
	}

	if value != "http://localhost/terraform/1.7.0/terraform_1.7.0_darwin_amd64.zip" {
		t.Error("Unexpected result, get :", value)
	}
}
