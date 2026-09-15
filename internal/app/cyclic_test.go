package app

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestMessageViewFromDescGoogleProtobufValue(t *testing.T) {
	md := (&structpb.Value{}).ProtoReflect().Descriptor()
	view, err := messageViewFromDesc(md, &cyclicDetector{})
	if err != nil {
		t.Fatalf("google.protobuf.Value should be usable (snip cycles), got: %v", err)
	}
	if view == nil || view.FullName != "google.protobuf.Value" {
		t.Fatalf("unexpected view: %+v", view)
	}
	if len(view.Fields) == 0 {
		t.Fatal("expected Value fields, got none")
	}
}

func TestMessageViewFromDescSelfReferentialSnips(t *testing.T) {
	node := mustSelfRefNode(t)
	view, err := messageViewFromDesc(node, &cyclicDetector{})
	if err != nil {
		t.Fatalf("self-referential message should snip at depth, not fail: %v", err)
	}
	if view == nil {
		t.Fatal("expected non-nil view")
	}

	depth := 0
	for cur := view; cur != nil; {
		depth++
		if depth > 20 {
			t.Fatal("cycle was not truncated")
		}
		if len(cur.Fields) == 0 {
			break
		}
		cur = cur.Fields[0].Message
	}
	// maxCyclicDepth successful Node visits, then one snipped empty leaf.
	if depth != maxCyclicDepth+1 {
		t.Fatalf("expected nesting depth %d (snip after %d), got %d", maxCyclicDepth+1, maxCyclicDepth, depth)
	}
}

func TestCyclicDetectorErrorsBeyondMaxDepth(t *testing.T) {
	md := (&structpb.Value{}).ProtoReflect().Descriptor()
	cd := &cyclicDetector{}
	for i := 0; i < maxCyclicDepth; i++ {
		if err := cd.detect(md); err != nil {
			t.Fatalf("detect at visit %d (of %d allowed): %v", i+1, maxCyclicDepth, err)
		}
	}
	err := cd.detect(md)
	if err == nil {
		t.Fatal("expected error once max cyclic depth exceeded")
	}
	if !strings.Contains(err.Error(), "cyclic") {
		t.Fatalf("expected cyclic error, got %v", err)
	}
	// Failed detect must not push path/graph (would desync pop).
	if len(cd.path) != maxCyclicDepth || len(cd.graph) != maxCyclicDepth {
		t.Fatalf("snip should not push; path=%d graph=%d want %d", len(cd.path), len(cd.graph), maxCyclicDepth)
	}
}

func TestCyclicDetectorSnipKeepsPathGraphInSync(t *testing.T) {
	node := mustSelfRefNode(t)
	cd := &cyclicDetector{}
	if _, err := messageViewFromDesc(node, cd); err != nil {
		t.Fatal(err)
	}
	if len(cd.path) != 0 || len(cd.graph) != 0 {
		t.Fatalf("after full walk path/graph should be empty; path=%v graph=%v", cd.path, cd.graph)
	}
}

func mustSelfRefNode(t *testing.T) protoreflect.MessageDescriptor {
	t.Helper()
	name := "node.proto"
	pkg := "test"
	syntax := "proto3"
	msgName := "Node"
	fieldName := "child"
	typeName := ".test.Node"
	num := int32(1)
	label := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	ftype := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE

	fd, err := protodesc.NewFile(&descriptorpb.FileDescriptorProto{
		Name:    &name,
		Package: &pkg,
		Syntax:  &syntax,
		MessageType: []*descriptorpb.DescriptorProto{{
			Name: &msgName,
			Field: []*descriptorpb.FieldDescriptorProto{{
				Name:     &fieldName,
				Number:   &num,
				Label:    &label,
				Type:     &ftype,
				TypeName: &typeName,
			}},
		}},
	}, nil)
	if err != nil {
		t.Fatalf("protodesc.NewFile: %v", err)
	}
	return fd.Messages().Get(0)
}
