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
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Request(callURL string, selector string, extracter func(*goquery.Selection) string) ([]string, error) {
	response, err := http.Get(callURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return extract(response.Body, selector, extracter)
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

func extract(reader io.Reader, selector string, extracter func(*goquery.Selection) string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var extracteds []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if extracted := extracter(s); extracted != "" {
			extracteds = append(extracteds, extracted)
		}
	})

	return extracteds, nil
}

func selectionTextExtractor(s *goquery.Selection) string {
	return strings.TrimSpace(s.Text())
}
