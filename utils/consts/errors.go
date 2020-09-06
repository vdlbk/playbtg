package consts

const (
	ErrCannotOpenFile = "fail to open file"
)

type DefaultError string

func (d DefaultError) Error() string {
	return string(d)
}
