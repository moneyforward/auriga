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

package model

// ReactionSkinTones is all of skin-tone names supported by Slack
// A reaction with skin-tone can be expressed by combining the following array elements with the reaction specified by the user
var ReactionSkinTones []string = []string{
	"", // this value is used to express reactions without skin-tone
	":skin-tone-2:",
	":skin-tone-3:",
	":skin-tone-4:",
	":skin-tone-5:",
	":skin-tone-6:",
}
