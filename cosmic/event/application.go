package event

type AppTick struct {
	baseEvent
}

func (e *AppTick) Type() Type {
	return TypeAppTick
}

func (e *AppTick) Category() Category {
	return CategoryApplication
}

func (e *AppTick) String() string {
	return "AppTickEvent"
}

type AppUpdate struct {
	baseEvent
}

func (e *AppUpdate) Type() Type {
	return TypeAppUpdate
}

func (e *AppUpdate) Category() Category {
	return CategoryApplication
}

func (e *AppUpdate) String() string {
	return "AppUpdateEvent"
}

type AppRender struct {
	baseEvent
}

func (e *AppRender) Type() Type {
	return TypeAppRender
}

func (e *AppRender) Category() Category {
	return CategoryApplication
}

func (e *AppRender) String() string {
	return "AppRenderEvent"
}
