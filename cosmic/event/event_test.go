package event

import "testing"

type testEvent struct {
	*baseEvent
	category Category
}

func (te testEvent) Type() Type {
	return TypeAppRender
}

func (te testEvent) Category() Category {
	return te.category
}

func (te testEvent) String() string {
	return "TestEvent"
}

func TestIsInCategory(t *testing.T) {
	e := testEvent{}

	allExceptNone := []Category{
		CategoryApplication,
		CategoryWindow,
		CategoryInput,
		CategoryKey,
		CategoryMouse,
		CategoryMouseButton,
	}
	all := append(allExceptNone, CategoryNone)

	for _, category := range all {
		if IsInCategory(e, category) {
			t.Error("should not match any category when set to 0")
		}
	}

	for _, set := range allExceptNone {
		e.category = set

		for _, category := range all {
			if category == set && !IsInCategory(e, category) {
				t.Errorf("expected true (%02x IsInCategory %02x)", set, category)
			}

			if category != set && IsInCategory(e, category) {
				t.Errorf("expected false (%02x IsInCategory %02x)", set, category)
			}
		}
	}

	e.category = 0xffff
	for _, category := range allExceptNone {
		if !IsInCategory(e, category) {
			t.Errorf("should match every category when set to all ones (IsInCategory %02x)", category)
		}
	}
}

func Test_baseEvent_IsHandled(t *testing.T) {
	handledEvent := baseEvent{handled: true}
	if !handledEvent.IsHandled() {
		t.Error("handled event should return true")
	}

	unhandledEvent := baseEvent{handled: false}
	if unhandledEvent.IsHandled() {
		t.Error("unhandled event should return false")
	}
}

func Test_baseEvent_SetHandled(t *testing.T) {
	testEvent := baseEvent{}

	if testEvent.handled {
		t.Error("new events should not be handled")
	}

	testEvent.SetHandled()
	if !testEvent.handled {
		t.Error("event.handled should be true after calling SetHandled")
	}
}
