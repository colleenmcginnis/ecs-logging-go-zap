// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package ecszap

import (
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

// EpochMicrosTimeEncoder encodes a given time in microseconds.
func EpochMicrosTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	micros := float64(t.UnixNano()) / float64(time.Microsecond)
	enc.AppendFloat64(micros)
}

// CallerEncoder is equivalent to zapcore.CallerEncoder, except that its UnmarshalText
// method uses FullCallerEncoder and ShortCallerEncoder from this package instead,
// in order to encode callers in the ECS format.
type CallerEncoder func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder)

// FullCallerEncoder serializes the file name, line and function from the caller
// in an ECS compliant way; serializing the full path of the file name
// using the underlying zapcore.EntryCaller.
func FullCallerEncoder(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	encodeCaller(&caller{c, true}, enc)
}

// ShortCallerEncoder serializes the file name, line and function from the caller
// in an ECS compliant way; removing everything except the final directory from the
// file name by calling the underlying zapcore.EntryCaller TrimmedPath().
func ShortCallerEncoder(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	encodeCaller(&caller{c, false}, enc)
}

// RFC3339TimeEncoder serializes a time.Time to an RFC3339-formatted string
// with millisecond precision.
//
// If enc supports AppendTimeLayout(t time.Time,layout string), it's used
// instead of appending a pre-formatted string value.
func RFC3339TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}
	// Use a custom RFC3339 layout with obligatory millisecond
	// precision rather than the second-resolution or optional
	// nano-second resolution that is provided in the time
	// package.
	const rfc3339millis = "2006-01-02T15:04:05.000Z07:00"
	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, rfc3339millis)
		return
	}

	enc.AppendString(t.Format(rfc3339millis))
}

// RFC3339UTCTimeEncoder serializes a time.Time to an RFC3339-formatted string
// with millisecond precision in UTC.
//
// If enc supports AppendTimeLayout(t time.Time,layout string), it's used
// instead of appending a pre-formatted string value.
//
// RFC3339UTCTimeEncoder is the default time encoder used by EncoderConfig.
func RFC3339UTCTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	RFC3339TimeEncoder(t.In(time.UTC), enc)
}

// UnmarshalText creates a CallerEncoder function,
// `full` is unmarshalled to FullCallerEncoder,
// defaults to ShortCallerEncoder,
func (e *CallerEncoder) UnmarshalText(text []byte) error {
	switch string(text) {
	case "full":
		*e = FullCallerEncoder
	default:
		*e = ShortCallerEncoder
	}
	return nil
}

func encodeCaller(c *caller, enc zapcore.PrimitiveArrayEncoder) {
	// this function can only be called internally so we have full control over it
	// and can ensure that enc is always of type zapcore.ArrayEncoder
	if e, ok := enc.(zapcore.ArrayEncoder); ok {
		e.AppendObject(c)
	}
}

type caller struct {
	zapcore.EntryCaller
	fullPath bool
}

func (c *caller) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var file string
	if c.fullPath {
		file = c.File
	} else {
		file = c.TrimmedPath()
		file = file[:strings.LastIndex(file, ":")]
	}
	enc.AddString("function", c.Function)
	enc.AddString("file.name", file)
	enc.AddInt("file.line", c.Line)
	return nil
}
