alias cdo="cd ~/src/qserv-operator"

# Rebuild and re-install qserv from scratch
alias rbo="cdo && ./build.sh -k && make uninstall && kubectl delete deployment -n qserv-operator-system --all && \
	./deploy.sh && kubectl delete pvc --all && kubectl apply -k manifests/base && cd -"
