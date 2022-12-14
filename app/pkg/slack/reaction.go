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

package slack

import (
	"regexp"
	"strings"
)

// IsReaction judges whether an argument is reaction string or not.
func IsReaction(reaction string) bool {
	return strings.HasPrefix(reaction, ":") && strings.HasSuffix(reaction, ":")
}

// ExtractReactionName extracts reaction name (reaction) from reaction string (:reaction:)
func ExtractReactionName(reaction string) string {
	return strings.ReplaceAll(reaction, ":", "")
}

// regReactionSkinTone is regexp which indicate skin-tone names supported by Slack
var regReactionSkinTone = regexp.MustCompile(`skin-tone-\d+`)

// RemoveSkinToneFromReaction uses regexp to remove the skin-tone string
func RemoveSkinToneFromReaction(reaction string) string {
	return regReactionSkinTone.ReplaceAllString(reaction, "")
}
