FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/qserv-operator \
    USER_UID=1001 \
    USER_NAME=qserv-operator

# install operator binary
COPY build/_output/bin/qserv-operator ${OPERATOR}

COPY build/bin /usr/local/bin
COPY configmap /configmap
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
