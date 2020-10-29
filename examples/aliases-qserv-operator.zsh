alias cdo="cd ~/src/qserv-operator"
alias cdi="cd ~/src/qserv-ingest"

# Re-install qserv from scratch
alias krq="cdo && kubectl delete qservs.qserv.lsst.org --all && kubectl delete pvc,pv --all && kubectl apply -k manifests/base"

# Rebuild qserv from scratch
alias rbo="cdo && ./build.sh -k && ./push-image.sh -d && \
           kubectl delete deployment -n qserv-operator-system --all && \
	   ./deploy.sh && krq"

# Relaunch ingest
alias ri="cdi && ./build-image.sh && ./job.sh init && ./job.sh ingest && ./job.sh publish && ./job.sh index-tables"
