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

/*
Package errors は、 error を操作する関数を提供する。

goのデフォルト error とともに、goのエラーで提供しないstack tracingのエラーも表示したいので、両方を出力できるようにしている。

"errors" + "github.com/pkg/errors"(スタックトレースのため) のwrapper。

stack情報をerrorに保存するためには、Wrap, Wrapfで errorにstack traceを含めるようにwrapして、

出力する側で %+v 形で表示する。
*/
package errors
