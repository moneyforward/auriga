/*
 * Copyright 2022 Money Forward, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package slice

// ToStringSet remove duplicated string element
// TODO: Update 1.18 and use generics
func ToStringSet(s []string) []string {
	var set []string
	m := map[string]bool{}
	for _, e := range s {
		if !m[e] {
			m[e] = true
			set = append(set, e)
		}
	}
	return set
}
