package issue22

type Bar int

func (yyy Bar) String() string {
	return "bar"
}

func (zzz Bar) String2() string {
	return "bar2"
}
