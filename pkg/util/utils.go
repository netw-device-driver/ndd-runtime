/*
Copyright 2021 Wim Henderickx.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

// BoolPtr return pointer to boolean
func BoolPtr(b bool) *bool { return &b }

// StringPtr return pointer to boolean
func StringPtr(s string) *string { return &s }

// intPtr return pointer to int
func intPtr(i int) *int { return &i }

// int32Ptr return pointer to int32
func int32Ptr(i int32) *int32 { return &i }

// Uint32Ptr return pointer to uint32
func Uint32Ptr(ui uint32) *uint32 { return &ui }
