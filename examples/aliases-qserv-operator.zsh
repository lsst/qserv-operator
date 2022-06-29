
QSERV_OPERATOR_SRC_DIR="~/src/qserv-operator"
QSERV_INGEST_SRC_DIR="~/src/qserv-ingest"

alias cdo="cd $QSERV_OPERATOR_SRC_DIR"
alias cdi="cd $QSERV_INGEST_SRC_DIR"

# Aliases for qserv-operator
############################

# Re-install qserv from scratch

alias delete-qserv="kubectl delete qservs.qserv.lsst.org --all && kubectl delete pvc -l app.kubernetes.io/managed-by=qserv-operator"
alias krq="delete-qserv && kubectl apply -k $QSERV_OPERATOR_SRC_DIR/manifests/base"
alias gkrq="delete-qserv && kubectl apply -k $QSERV_OPERATOR_SRC_DIR/manifests/gke-qserv-dev"

# Rebuild qserv from scratch
alias rbo="cdo && ./build.sh -k && ./push-image.sh -k && \
           kubectl delete deployment,pod -n qserv-operator-system --all && \
           ./deploy.sh && krq"

# Aliases for qserv-ingest
##########################

# Restart ingest
alias arsub="cdi && ./build.sh && ./push-image.sh && ./argo-submit.sh"
alias arw="argo watch @latest"
alias arg="argo get @latest"
# Delete all previous ingests
alias ardel="argo delete --all && kubectl delete job -l app=qserv,tier=ingest"
# Delete then restart ingest
alias arrestart="ardel && cdi && ./example/delete_database.sh qservTest_case01_qserv CHANGEME && arsub"
