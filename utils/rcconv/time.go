// Copyright (c) 2024 Focela Technologies. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rcconv

import (
	"time"

	"github.com/focela/ratcatcher/internal/utils"
	"github.com/focela/ratcatcher/os/rctime"
)

// Time converts `any` to time.Time.
func Time(any interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v
		}
	}
	if t := RCTime(any, format...); t != nil {
		return t.Time
	}
	return time.Time{}
}

// Duration converts `any` to time.Duration.
// If `any` is string, then it uses time.ParseDuration to convert it.
// If `any` is numeric, then it converts `any` as nanoseconds.
func Duration(any interface{}) time.Duration {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v
	}
	s := String(any)
	if !utils.IsNumeric(s) {
		d, _ := rctime.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(any))
}

// RCTime converts `any` to *rctime.Time.
// The parameter `format` can be used to specify the format of `any`.
// It returns the converted value that matched the first format of the formats slice.
// If no `format` given, it converts `any` using rctime.NewFromTimeStamp if `any` is numeric,
// or using rctime.StrToTime if `any` is string.
func RCTime(any interface{}, format ...string) *rctime.Time {
	if any == nil {
		return nil
	}
	if v, ok := any.(iRCTime); ok {
		return v.RCTime(format...)
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(*rctime.Time); ok {
			return v
		}
		if t, ok := any.(time.Time); ok {
			return rctime.New(t)
		}
		if t, ok := any.(*time.Time); ok {
			return rctime.New(t)
		}
	}
	s := String(any)
	if len(s) == 0 {
		return rctime.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		for _, item := range format {
			t, err := rctime.StrToTimeFormat(s, item)
			if t != nil && err == nil {
				return t
			}
		}
		return nil
	}
	if utils.IsNumeric(s) {
		return rctime.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := rctime.StrToTime(s)
		return t
	}
}