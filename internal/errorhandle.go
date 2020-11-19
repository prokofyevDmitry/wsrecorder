package errhdl

func PanicIf(e error) {
	if e != nil {
		panic(e)
	}
}

func LogIf(e error) {
	if e != nil {
		panic(e)
	}
}
