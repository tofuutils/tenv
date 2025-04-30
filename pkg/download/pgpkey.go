/*
 *
 * Copyright 2025 tofuutils authors.
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

package download

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	DefaultHashicorpPGPKeyURL = "https://www.hashicorp.com/.well-known/pgp-key.txt"
)

// GetPGPKey retrieves the PGP key from either a local file or URL.
// If keyPath starts with "http://" or "https://", it will be treated as a URL.
// Otherwise, it will be treated as a local file path.
// If keyPath is empty, it will use the default HashiCorp PGP key URL.
func GetPGPKey(ctx context.Context, keyPath string, display func(string)) ([]byte, error) {
	if keyPath == "" {
		return Bytes(ctx, DefaultHashicorpPGPKeyURL, display, func(resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to download PGP key: HTTP %d", resp.StatusCode)
			}
			return nil
		})
	}

	if strings.HasPrefix(keyPath, "http://") || strings.HasPrefix(keyPath, "https://") {
		return Bytes(ctx, keyPath, display, func(resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("failed to download PGP key: HTTP %d", resp.StatusCode)
			}
			return nil
		})
	}

	// Try to read from local file
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PGP key from %s: %w", keyPath, err)
	}
	return data, nil
}
