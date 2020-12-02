packge merr

type MyError struct {
	File string
	Line int
	Msg string
}

func (e *MyError) Error() string {
	return fmt.Sprinf("%s: %d: %s", e.File, e.Line, e.Msg)
}

func NewTest() error {
	return &MyError{"sever.go", 10, "someting error..."}
}

func errSwitch(err error) {
	switch errtype := err.(type) {
	case nil:
		// do nothing
	case *MyError:
		fmt.Println(err.Error(), errtype.Line)
	default:
		// unkonw error
	}
}