/*
Copyright 2017 Google Inc.

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

package parser

import (
	"bytes"
	"fmt"
)

// TrackedBuffer is used to rebuild a query from the ast.
type TrackedBuffer struct {
	*bytes.Buffer
}

// NewTrackedBuffer creates a new TrackedBuffer.
func NewTrackedBuffer() *TrackedBuffer {
	return &TrackedBuffer{
		Buffer: new(bytes.Buffer),
	}
}

// WriteNode function, initiates the writing of a single SQLNode tree by passing
// through to Myprintf with a default format string
func (buf *TrackedBuffer) WriteNode(node SQLNode) *TrackedBuffer {
	buf.Myprintf("%v", node)
	return buf
}

// Myprintf mimics fmt.Fprintf(buf, ...), but limited to Node(%v),
// Node.Value(%s) and string(%s). It also allows a %a for a value argument, in
// which case it adds tracking info for future substitutions.
//
// The name must be something other than the usual Printf() to avoid "go vet"
// warnings due to our custom format specifiers.
func (buf *TrackedBuffer) Myprintf(format string, values ...interface{}) {
	end := len(format)
	fieldnum := 0
	for i := 0; i < end; {
		lasti := i
		for i < end && format[i] != '%' {
			i++
		}
		if i > lasti {
			buf.WriteString(format[lasti:i])
		}
		if i >= end {
			break
		}
		i++ // '%'
		switch format[i] {
		case 'c':
			switch v := values[fieldnum].(type) {
			case byte:
				buf.WriteByte(v)
			case rune:
				buf.WriteRune(v)
			default:
				panic(fmt.Sprintf("unexpected TrackedBuffer type %T", v))
			}
		case 's':
			switch v := values[fieldnum].(type) {
			case []byte:
				buf.Write(v)
			case string:
				buf.WriteString(v)
			default:
				panic(fmt.Sprintf("unexpected TrackedBuffer type %T", v))
			}
		case 'v':
			node := values[fieldnum].(SQLNode)
			node.Format(buf)
		case 'a':
			buf.WriteString(values[fieldnum].(string))
		default:
			panic("unexpected")
		}
		fieldnum++
		i++
	}
}
