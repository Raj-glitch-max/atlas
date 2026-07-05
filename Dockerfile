# Engineering Foundation — dev-toolbox container (stack-agnostic).
# This is NOT a product container. Product containers are stack-specific and are
# deferred until a stack is selected post-C4-spike (see CONTRIBUTING.md §5).
#
# Build: make devshell   (or: docker build -t eng-toolbox . )
# Run:   make devshell
#
# Scope: a reproducible contributor shell with the Sprint-0 toolchain.

FROM debian:12-slim

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=C.UTF-8 \
    LC_ALL=C.UTF-8

RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates curl git jq make shellcheck yamllint \
      python3 python3-pip python3-venv \
      less openssh-client \
    && pip3 install --no-cache-dir --break-system-packages pre-commit \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /ws
CMD ["bash"]
