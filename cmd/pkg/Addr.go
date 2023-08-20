package pkg

type Addr struct {
	Net     string
	Address string
}

func (a Addr) Network() string {
	return a.Net
}
func (a Addr) String() string {
	return a.Address
}
