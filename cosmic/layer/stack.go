package layer

type Stack struct {
	layers   []Layer
	overlays []Layer
}

// ForEachAscending executes the given function for all layers and overlays in
// ascending order (starting from the bottom of the stack).
func (ls *Stack) ForEachAscending(fn func(l Layer)) {
	for _, layer := range ls.layers {
		fn(layer)
	}

	for _, overlay := range ls.overlays {
		fn(overlay)
	}
}

// ForEachDescending executes the given function for all layers and overlays in
// descending order (starting from the top of the stack).
func (ls *Stack) ForEachDescending(fn func(l Layer)) {
	for i := range ls.overlays {
		fn(ls.overlays[len(ls.overlays)-1-i])
	}

	for i := range ls.layers {
		fn(ls.layers[len(ls.layers)-1-i])
	}
}

func (ls *Stack) Push(l Layer) {
	l.OnAttach()
	ls.layers = append(ls.layers, l)
}

func (ls *Stack) Pop(l Layer) (popped Layer) {
	l.OnDetach()
	popped, ls.layers = ls.layers[len(ls.layers)-1], ls.layers[:len(ls.layers)-1]

	return
}

func (ls *Stack) PushOverlay(l Layer) {
	l.OnAttach()
	ls.overlays = append(ls.overlays, l)
}

func (ls *Stack) PopOverlay(l Layer) (popped Layer) {
	l.OnDetach()
	popped, ls.overlays = ls.overlays[len(ls.overlays)-1], ls.overlays[:len(ls.overlays)-1]

	return
}
