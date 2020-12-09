package errorf

import "encoding/json"

type CErr struct {
	Msg  string
	Args interface{}
}

func (e CErr) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

func New(msg string, args ...interface{}) CErrF {
	return CErrF{
		CErr{
			Msg:  msg,
			Args: args,
		},
	}
}

type CErrF struct {
	CErr
}

type CErrFs []CErrF

func (errs CErrFs) Add(newError CErrF) CErrFs {
	return append(errs, newError)
}

func (errs CErrFs) Error() string {
	return ""
}
func (errs CErrFs) Get() CErrFs {
	return errs
}
