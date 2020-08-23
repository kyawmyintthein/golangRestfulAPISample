package newrelicx

import "github.com/kyawmyintthein/orange-contrib/errorx"

type NotAvailable struct {
	*errorx.ErrorX
}

func NotAvailableError() *NotAvailable {
	return &NotAvailable{
		errorx.NewErrorX("[%s] new-relic tracer is not avaliable", PackageName),
	}
}
