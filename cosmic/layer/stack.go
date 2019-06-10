package layer

// Stack represents a collection of layers. It differentiates between layers
// and overlays (both represented by a Layer), where overlays sit 'on top' of
// layers.
type Stack struct {
	layers   []Layer
	overlays []Layer
}

func (ls *Stack) Push(l Layer) {
	l.OnAttach()
	ls.layers = append(ls.layers, l)
}

func (ls *Stack) Pop() (popped Layer) {
	popped, ls.layers = ls.layers[len(ls.layers)-1], ls.layers[:len(ls.layers)-1]
	popped.OnDetach()

	return
}

func (ls *Stack) PushOverlay(l Layer) {
	l.OnAttach()
	ls.overlays = append(ls.overlays, l)
}

func (ls *Stack) PopOverlay() (popped Layer) {
	popped, ls.overlays = ls.overlays[len(ls.overlays)-1], ls.overlays[:len(ls.overlays)-1]
	popped.OnDetach()

	return
}

// Bottom provides a new StackItem which can be used to iterate from the bottom
// of the Stack.
func (ls *Stack) Bottom() *StackItem {
	if len(ls.layers) > 0 {
		return &StackItem{
			stack: ls,
			index: -1,
		}
	}

	if len(ls.overlays) > 0 {
		return &StackItem{
			stack: ls,
			index: len(ls.layers) - 1,
		}
	}

	panic("cannot get the bottom layer of an empty layer stack")
}

// Top provides a new StackItem which can be used to iterate from the top of
// the Stack.
func (ls *Stack) Top() *StackItem {
	if len(ls.overlays) > 0 {
		return &StackItem{
			stack: ls,
			index: len(ls.layers) + len(ls.overlays),
		}
	}

	if len(ls.layers) > 0 {
		return &StackItem{
			stack: ls,
			index: len(ls.layers),
		}
	}

	panic("cannot get the top layer of an empty layer stack")
}

// StackItem allows callers to iterate through a Stack. Please note that this
// is NOT thread-safe, and will break if the underlying Stack is mutated while
// being iterated over.
type StackItem struct {
	stack *Stack
	index int
}

// Get retrieves the currently selected Layer in the Stack. When Get is
// called after Next or Prev return false, it returns nil.
func (item *StackItem) Get() Layer {
	if item.index >= len(item.stack.layers) {
		return item.stack.overlays[item.index-len(item.stack.layers)]
	} else if item.index >= 0 {
		return item.stack.layers[item.index]
	} else {
		return nil
	}
}

// Next iterates through the underlying Stack, selecting the next available
// Layer. It returns whether this operation succeeded.
func (item *StackItem) Next() bool {
	item.index++
	return item.index < len(item.stack.layers)+len(item.stack.overlays)
}

// Prev iterates through the underlying Stack, selecting the previous available
// Layer. It returns whether this operation succeeded.
func (item *StackItem) Prev() bool {
	item.index--
	return item.index >= 0
}
