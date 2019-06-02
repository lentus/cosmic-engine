package layer

import "github.com/lentus/cosmic-engine/cosmic/event"

type Layer interface {
	OnAttach()
	OnDetach()
	OnUpdate()
	OnEvent(e event.Event)
}
