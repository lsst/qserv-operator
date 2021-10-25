
QSERV_OPERATOR_SRC_DIR="~/src/qserv-operator"
QSERV_INGEST_SRC_DIR="~/src/qserv-ingest"

alias cdo="cd $QSERV_OPERATOR_SRC_DIR"
alias cdi="cd $QSERV_INGEST_SRC_DIR"


# Re-install qserv from scratch

alias delete-qserv="kubectl delete qservs.qserv.lsst.org --all && kubectl delete pvc -l app.kubernetes.io/managed-by=qserv-operator"
alias krq="delete-qserv && kubectl apply -k $QSERV_OPERATOR_SRC_DIR/manifests/base"
alias gkrq="delete-qserv && kubectl apply -k $QSERV_OPERATOR_SRC_DIR/manifests/gke-qserv-dev"

# Rebuild qserv from scratch
alias rbo="cdo && ./build.sh -k && ./push-image.sh -d && \
           kubectl delete deployment -n qserv-operator-system --all && \
	   ./deploy.sh && krq"

# Relaunch ingest
alias ri="cdi && ./build-image.sh && ./argo-submit.sh"
alias aw="argo watch @latest"
alias ag="argo get @latest"
