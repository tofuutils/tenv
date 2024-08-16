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

package htmlquery

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/tofuutils/tenv/v3/pkg/download"
)

func Request(callURL string, selector string, extractor func(*goquery.Selection) string, ro ...download.RequestOption) ([]string, error) {
	data, err := download.Bytes(callURL, download.NoDisplay, ro...)
	if err != nil {
		return nil, err
	}

	return extractList(data, selector, extractor)
}

func SelectionExtractor(part string) func(*goquery.Selection) string {
	if part == "#text" {
		return selectionTextExtractor
	}

	return func(s *goquery.Selection) string {
		attr, _ := s.Attr(part)

		return strings.TrimSpace(attr)
	}
}

func extractList(data []byte, selector string, extractor func(*goquery.Selection) string) ([]string, error) {
	dataReader := bytes.NewReader(data)
	doc, err := goquery.NewDocumentFromReader(dataReader)
	if err != nil {
		return nil, err
	}

	var extracteds []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if extracted := extractor(s); extracted != "" {
			extracteds = append(extracteds, extracted)
		}
	})

	return extracteds, nil
}

func selectionTextExtractor(s *goquery.Selection) string {
	return strings.TrimSpace(s.Text())
}
