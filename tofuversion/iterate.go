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

package tofuversion

import "sync"

func iterate[T any](values []T, reverseOrder bool) (<-chan T, func()) {
	valueChan := make(chan T)
	doneChan := make(chan struct{})
	if reverseOrder {
		go innerReverseIterate(values, valueChan, doneChan)
	} else {
		go innerIterate(values, valueChan, doneChan)
	}
	var once sync.Once
	return valueChan, func() {
		once.Do(func() {
			close(doneChan)
		})
	}
}

func innerIterate[T any](values []T, valueSender chan<- T, doneReceiver <-chan struct{}) {
ForLoop:
	for _, value := range values {
		select {
		case valueSender <- value:
		case <-doneReceiver:
			break ForLoop
		}
	}
	close(valueSender)
}

func innerReverseIterate[T any](values []T, valueSender chan<- T, doneReceiver <-chan struct{}) {
ForLoop:
	for i := len(values) - 1; i >= 0; i-- {
		select {
		case valueSender <- values[i]:
		case <-doneReceiver:
			break ForLoop
		}
	}
	close(valueSender)
}
