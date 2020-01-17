package v1alpha1

const (
	defaultWorkerNumber = 3
	defaultQservImage   = "qserv/qserv:11a6001"
)

// SetDefaults sets Qserv default fields
func (r *Qserv) SetDefaults() {

	if r.Spec.Worker.Replicas == 0 {
		r.Spec.Worker.Replicas = defaultWorkerNumber
	}

	if len(r.Spec.Worker.Image) == 0 {
		r.Spec.Worker.Image = defaultQservImage
	}

	if len(r.Spec.Czar.Image) == 0 {
		r.Spec.Czar.Image = defaultQservImage
	}
}
