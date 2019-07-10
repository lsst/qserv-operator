package v1alpha1

const (
	defaultWorkerNumber = 3
	defaultQservImage   = "qserv/qserv:11a6001"
)

// SetDefaults sets Redis field defaults
func (r *Qserv) SetDefaults() {

	if r.Spec.Worker.Replicas == 0 {
		r.Spec.Worker.Replicas = defaultWorkerNumber
	}

	if len(r.Spec.Worker.Image) == 0 {
		r.Spec.Worker.Image = defaultQservImage
	}
}
