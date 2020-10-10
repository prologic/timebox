package event

// TypedApplier calls the Applier appropriate for an Event's Type
type TypedApplier map[Type]Applier

// Applier dispatches to an underlying registered Applier
func (t TypedApplier) Applier() Applier {
	return func(e *Event) {
		if apply, ok := t[e.Type]; ok {
			apply(e)
		}
	}
}

// Combine with other instances, yielding a new instance
func (t TypedApplier) Combine(appliers ...TypedApplier) TypedApplier {
	combined := append([]TypedApplier{t}, appliers...)
	return TypedAppliers(combined...)
}

// TypedAppliers combines a TypedApplier set, performing the fan
// out to multiple appliers for the same type
func TypedAppliers(appliers ...TypedApplier) TypedApplier {
	res := TypedApplier{}
	for _, typed := range appliers {
		for eventType, applier := range typed {
			if originalApplier, ok := res[eventType]; ok {
				res[eventType] = func(e *Event) {
					originalApplier(e)
					applier(e)
				}
			} else {
				res[eventType] = applier
			}
		}
	}
	return res
}
