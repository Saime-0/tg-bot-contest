package usageErr

type UsageErr struct {
	Err   error
	Usage string
}

func (u *UsageErr) Error() string {
	return u.Err.Error()
}
