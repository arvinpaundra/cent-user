package trait

type Removeable struct {
	val bool
}

func (r *Removeable) IsMarkedToBeRemoved() bool {
	return r.val
}

func (r *Removeable) MarkToBeRemoved() {
	r.val = true
}

func (r *Removeable) UnmarkToBeRemoved() {
	r.val = false
}
