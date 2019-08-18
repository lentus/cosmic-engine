package layer

import (
	"github.com/lentus/cosmic-engine/cosmic/event"
	"testing"
)

// Provides a type implementing the Layer interface for use in the below tests.
type TestLayer struct {
	name string
}

func (el *TestLayer) OnAttach() {
}

func (el *TestLayer) OnDetach() {
}

func (el *TestLayer) OnUpdate() {
}

func (el *TestLayer) OnEvent(e event.Event) {
}

func TestStack_Push(t *testing.T) {
	stack := Stack{}
	stack.Push(&TestLayer{})

	if len(stack.layers) != 1 {
		t.Errorf("expected stack to contain 1 layer, got %d", len(stack.layers))
	}
}

func TestStack_Pop(t *testing.T) {
	stack := Stack{
		layers: []Layer{&TestLayer{name: "layer 1"}, &TestLayer{name: "layer 2"}},
	}

	popped := stack.Pop()
	layerName := popped.(*TestLayer).name

	if layerName != "layer 2" {
		t.Errorf("expected popped layer to be layer 2, got %s", layerName)
	}
	if len(stack.layers) != 1 {
		t.Errorf("expected stack to contain 1 layer, got %d", len(stack.layers))
	}
}

func TestStack_PushOverlay(t *testing.T) {
	stack := Stack{}
	stack.PushOverlay(&TestLayer{})

	if len(stack.overlays) != 1 {
		t.Errorf("expected stack to contain 1 overlay, got %d", len(stack.overlays))
	}
}

func TestStack_PopOverlay(t *testing.T) {
	stack := Stack{
		overlays: []Layer{&TestLayer{name: "overlay 1"}, &TestLayer{name: "overlay 2"}},
	}
	popped := stack.PopOverlay()
	overlayName := popped.(*TestLayer).name

	if overlayName != "overlay 2" {
		t.Errorf("expected popped overlay to be overlay 2, got %s", overlayName)
	}
	if len(stack.overlays) != 1 {
		t.Errorf("expected stack to contain 1 overlay, got %d", len(stack.overlays))
	}
}

func TestStack_Bottom_with_layers(t *testing.T) {
	stack := Stack{
		layers:   []Layer{&TestLayer{name: "layer 1"}, &TestLayer{name: "layer 2"}},
		overlays: []Layer{&TestLayer{name: "overlay 1"}, &TestLayer{name: "overlay 2"}},
	}

	bottom := stack.Bottom()
	layerName := bottom.Get().(*TestLayer).name

	if layerName != "layer 1" {
		t.Errorf("expected bottom stack item to be layer 1, got %s", layerName)
	}
}

func TestStack_Bottom_without_layers(t *testing.T) {
	stack := Stack{
		overlays: []Layer{&TestLayer{name: "overlay 1"}, &TestLayer{name: "overlay 2"}},
	}

	bottom := stack.Bottom()
	overlayName := bottom.Get().(*TestLayer).name

	if overlayName != "overlay 1" {
		t.Errorf("expected bottom stack item to be overlay 1, got %s", overlayName)
	}
}

func TestStack_Top_with_overlays(t *testing.T) {
	stack := Stack{
		layers:   []Layer{&TestLayer{name: "layer 1"}, &TestLayer{name: "layer 2"}},
		overlays: []Layer{&TestLayer{name: "overlay 1"}, &TestLayer{name: "overlay 2"}},
	}

	top := stack.Top()
	layerName := top.Get().(*TestLayer).name

	if layerName != "overlay 2" {
		t.Errorf("expected top stack item to be overlay 2, got %s", layerName)
	}
}

func TestStack_Top_without_overlays(t *testing.T) {
	stack := Stack{
		layers: []Layer{&TestLayer{name: "layer 1"}, &TestLayer{name: "layer 2"}},
	}

	top := stack.Top()
	layerName := top.Get().(*TestLayer).name

	if layerName != "layer 2" {
		t.Errorf("expected top stack item to be layer 2, got %s", layerName)
	}
}

func TestStackItem_Get(t *testing.T) {
	stack := Stack{
		layers:   []Layer{&TestLayer{name: "layer 1"}, &TestLayer{name: "layer 2"}},
		overlays: []Layer{&TestLayer{name: "overlay 1"}, &TestLayer{name: "overlay 2"}},
	}

	layerItem := StackItem{
		stack: &stack,
		index: 1,
	}

	layerName := layerItem.Get().(*TestLayer).name
	if layerName != "layer 2" {
		t.Errorf("expected stack item with index 1 to be layer 2, got %s", layerName)
	}

	overlayItem := StackItem{
		stack: &stack,
		index: 2,
	}

	overlayName := overlayItem.Get().(*TestLayer).name
	if overlayName != "overlay 1" {
		t.Errorf("expected stack item with index 2 to be overlay 1, got %s", overlayName)
	}
}

func TestStackItem_Next_success(t *testing.T) {
	stack := Stack{
		layers:   []Layer{&TestLayer{name: "layer"}},
		overlays: []Layer{&TestLayer{name: "overlay"}},
	}

	layerItem := StackItem{
		stack: &stack,
		index: 0,
	}

	if success := layerItem.Next(); !success {
		t.Error("expected Next() to succeed")
	}
	if layerItem.index != 1 {
		t.Errorf("expected stack index to be 1, got %d", layerItem.index)
	}
}

func TestStackItem_Next_failure(t *testing.T) {
	stack := Stack{
		layers: []Layer{&TestLayer{name: "layer"}},
	}

	layerItem := StackItem{
		stack: &stack,
		index: 0,
	}

	if success := layerItem.Next(); success {
		t.Error("expected Next() to fail")
	}
}

func TestStackItem_Prev_success(t *testing.T) {
	stack := Stack{
		layers:   []Layer{&TestLayer{name: "layer"}},
		overlays: []Layer{&TestLayer{name: "overlay"}},
	}

	layerItem := StackItem{
		stack: &stack,
		index: 1,
	}

	if success := layerItem.Prev(); !success {
		t.Error("expected Prev() to succeed")
	}
	if layerItem.index != 0 {
		t.Errorf("expected stack index to be 0, got %d", layerItem.index)
	}
}

func TestStackItem_Prev_failure(t *testing.T) {
	stack := Stack{
		layers: []Layer{&TestLayer{name: "layer"}},
	}

	layerItem := StackItem{
		stack: &stack,
		index: 0,
	}

	if success := layerItem.Prev(); success {
		t.Error("expected Prev() to fail")
	}
}
