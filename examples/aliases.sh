alias cdo="cd ~/src/qserv-operator"

# Rebuild and re-install qserv from scratch
alias rbo="cdo && ./build.sh && make uninstall && ./deploy.sh && kubectl delete pvc --all && kubectl apply -k base && cd -"
