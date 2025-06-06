package trait

type Updateable struct {
	val bool
}

func (u *Updateable) IsMarkedToBeUpdated() bool {
	return u.val
}

func (u *Updateable) MarkToBeUpdated() {
	u.val = true
}

func (u *Updateable) UnmarkToBeUpdated() {
	u.val = false
}
