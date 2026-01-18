package kube

var (
  LookupServer func(svc string) (*Server, error)
  LookupNet    func(svc string) (*Net, error)
  LookupIP     func(svc string) (string, error)
  LookupEPTs   func(svc string) (ips []string, err error) // get all ips of a service
)

func lookupServer(svc string) (*Server, error)
func lookupNet(svc string) (*Net, error) //
func lookupIP(svc string) (string, error)
func lookupEPTs(svc string) (ips []string, err error) // get all ips of a service
