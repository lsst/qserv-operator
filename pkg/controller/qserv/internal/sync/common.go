package sync

var controllerLabels = map[string]string{
	"app.kubernetes.io/managed-by": "qserv-operator",
}

var noFunc = func() error {
	return nil
}
