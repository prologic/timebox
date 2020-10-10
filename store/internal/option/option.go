package option

type (
	// Target is the value to which an option would be applied
	Target interface{}

	// Option applies an option to a configurable Target
	Option func(Target) error
)

// Apply applies Options to a Target
func Apply(target Target, options ...Option) error {
	for _, o := range options {
		if err := o(target); err != nil {
			return err
		}
	}
	return nil
}
