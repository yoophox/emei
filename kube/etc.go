package kube

func lookupServerInEtc(svc string) (*Server, error)
func lookupNetInEtc(svc string) (*Net, error) //
func lookupIPInEtc(svc string) (string, error)
func lookupEPTsInEtc(svc string) (ips []string, err error) // get all ips of a service
