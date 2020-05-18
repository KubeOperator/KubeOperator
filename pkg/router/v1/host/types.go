package host

type Host struct{}

type ListHostRequest struct {
	name string
}

type ListHostResponse struct {
	item Host
}
