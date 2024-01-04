# You can pick any Debian/Ubuntu-based image. 😊
FROM mcr.microsoft.com/devcontainers/go:1.21-bullseye

# Options for setup script
ARG INSTALL_ZSH="true"
ARG UPGRADE_PACKAGES="false"
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Install needed packages and setup non-root user. Use a separate RUN statement to add your own dependencies.
COPY library-scripts/*.sh /tmp/library-scripts/
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && /bin/bash /tmp/library-scripts/common-debian.sh "${INSTALL_ZSH}" "${USERNAME}" "${USER_UID}" "${USER_GID}" "${UPGRADE_PACKAGES}" "true" "true" \
    && apt-get install -y g++ \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/* /tmp/library-scripts

# Install Bazel
ARG BAZELISK_VERSION=v1.10.1
ARG BAZELISK_DOWNLOAD_SHA=dev-mode
RUN curl -fSsL -o /usr/local/bin/bazel https://github.com/bazelbuild/bazelisk/releases/download/${BAZELISK_VERSION}/bazelisk-linux-amd64 \
    && ([ "${BAZELISK_DOWNLOAD_SHA}" = "dev-mode" ] || echo "${BAZELISK_DOWNLOAD_SHA} */usr/local/bin/bazel" | sha256sum --check - ) \
    && chmod 0755 /usr/local/bin/bazel

# Install Buildifier
ARG BUILDIFIER_VERSION=4.2.4
ARG BUILDIFIER_DOWNLOAD_SHA=dev-mode
RUN curl -fSsL -o /usr/local/bin/buildifier https://github.com/bazelbuild/buildtools/releases/download/${BUILDIFIER_VERSION}/buildifier-linux-amd64 \
    && ([ "${BUILDIFIER_DOWNLOAD_SHA}" = "dev-mode" ] || echo "${BUILDIFIER_DOWNLOAD_SHA} */usr/local/bin/buildifier" | sha256sum --check - ) \
    && chmod 0755 /usr/local/bin/buildifier

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends python3-pip
