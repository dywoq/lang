// Copyright 2025 dywoq
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package debug

import (
	"io"
	"log"
	"sync/atomic"
)

type Debug struct {
	output io.Writer
	busy   *atomic.Bool
	logger *log.Logger
}

// New creates new [Debug] structure and returns a pointer of it.
func New(output io.Writer, busy *atomic.Bool) *Debug {
	l := log.New(output, "", 0)
	return &Debug{output: output, busy: busy, logger: l}
}

// SetOutput sets the output, but panics if it's busy.
func (d *Debug) SetOutput(output io.Writer) {
	if d.busy.Load() {
		panic("can't modify debug parameters when it's busy")
	}
	d.output = output
	d.logger.SetOutput(output)
}

// Output returns the set output instance.
func (d *Debug) Output() io.Writer {
	return d.output
}

func (d *Debug) Println(v ...any) {
	d.logger.Println(v...)
}

func (d *Debug) Printf(fmt string, v ...any) {
	d.logger.Printf(fmt, v...)
}

func (d *Debug) SetFlags(flags int) {
	d.logger.SetFlags(flags)
}
