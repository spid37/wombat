package app

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// maxCyclicDepth is how many times the same message type may appear on the
// descriptor walk path before we snip. google.protobuf.Value / ListValue are
// cyclic; real payloads rarely nest deeper than a handful of levels.
const maxCyclicDepth = 7

type cyclicDetector struct {
	path  []protoreflect.FullName
	graph []string
}

func (d *cyclicDetector) detect(md protoreflect.MessageDescriptor) error {
	n, f := md.Name(), md.FullName()
	count := 0
	for _, p := range d.path {
		if p == f {
			count++
		}
	}
	if count >= maxCyclicDepth {
		trail := strings.Join(d.graph, " → ")
		if trail != "" {
			trail += " → "
		}
		return fmt.Errorf("unable to parse proto descriptors: cyclic data detected: %s%s", trail, n)
	}
	// Only push after the depth check so a snip does not desync graph/path.
	d.path = append(d.path, f)
	d.graph = append(d.graph, string(n))
	return nil
}

func (d *cyclicDetector) pop() {
	if n := len(d.path); n > 0 {
		d.path = d.path[:n-1]
	}
	if n := len(d.graph); n > 0 {
		d.graph = d.graph[:n-1]
	}
}
