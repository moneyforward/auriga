/*
 * Copyright 2023 Money Forward, Inc.
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

import (
	"reflect"
	"testing"
)

func TestSplitStringSliceInChunks(t *testing.T) {
	type args struct {
		s         []string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "OK: all elements(slices) are the same length (chunk size)",
			args: args{
				s:         []string{"a", "b", "c", "d"},
				chunkSize: 2,
			},
			want: [][]string{
				{"a", "b"},
				{"c", "d"},
			},
		},
		{
			name: "OK: all elements(slices) are NOT the same length (chunk size)",
			args: args{
				s:         []string{"a", "b", "c", "d", "e"},
				chunkSize: 2,
			},
			want: [][]string{
				{"a", "b"},
				{"c", "d"},
				{"e"},
			},
		},
		{
			name: "OK: slices length is chunk size",
			args: args{
				s:         []string{"a", "b", "c", "d"},
				chunkSize: 4,
			},
			want: [][]string{
				{"a", "b", "c", "d"},
			},
		},
		{
			name: "OK: chunk size is bigger than slices length",
			args: args{
				s:         []string{"a", "b", "c", "d"},
				chunkSize: 5,
			},
			want: [][]string{
				{"a", "b", "c", "d"},
			},
		},
		{
			name: "OK: chunk size is less than 1",
			args: args{
				s:         []string{"a", "b", "c", "d"},
				chunkSize: 0,
			},
			want: [][]string{
				{"a", "b", "c", "d"},
			},
		},
		{
			name: "OK: s is empty and chunkSize is 0",
			args: args{
				s:         []string{},
				chunkSize: 0,
			},
			want: [][]string{{}},
		},
		{
			name: "OK: s is empty and chunkSize is more than 0",
			args: args{
				s:         []string{},
				chunkSize: 1,
			},
			want: [][]string{{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitStringSliceInChunks(tt.args.s, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitStringSliceInChunks() = %v, want %v", got, tt.want)
			}
		})
	}
}
