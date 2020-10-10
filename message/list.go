package message

// List represents an array of Message pointers
type List []*Message

// EmptyList is the empty Message list
var EmptyList = List{}

// HandleWith pipes this List's Messages into the designated
// Handler. Will short-circuit upon the first encountered error
// and return a list starting with the event in error
func (l List) HandleWith(handle Handler) (List, error) {
	for i, e := range l {
		if err := handle(e); err != nil {
			return l[i:], err
		}
	}
	return EmptyList, nil
}
