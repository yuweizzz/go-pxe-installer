#!/usr/bin/env bash

WORK_DIR=$(pwd)/help

# debian 12

# amd64
if [[ ! -f "${WORK_DIR}/images/debian-bookworm-amd64/linux" ]]; then
    curl -s \
    --create-dirs \
    -o ${WORK_DIR}/images/debian-bookworm-amd64/linux \
    https://deb.debian.org/debian/dists/bookworm/main/installer-amd64/current/images/netboot/debian-installer/amd64/linux
fi

if [[ ! -f "${WORK_DIR}/images/debian-bookworm-amd64/initrd.gz" ]]; then
    curl -s \
    --create-dirs \
    -o ${WORK_DIR}/images/debian-bookworm-amd64/initrd.gz \
    https://deb.debian.org/debian/dists/bookworm/main/installer-amd64/current/images/netboot/debian-installer/amd64/initrd.gz
fi

# arm64
if [[ ! -f "${WORK_DIR}/images/debian-bookworm-arm64/linux" ]]; then
    curl -s \
    --create-dirs \
    -o ${WORK_DIR}/images/debian-bookworm-arm64/linux \
    https://deb.debian.org/debian/dists/bookworm/main/installer-arm64/current/images/netboot/debian-installer/arm64/linux
fi

if [[ ! -f "${WORK_DIR}/images/debian-bookworm-arm64/initrd.gz" ]]; then
    curl -s \
    --create-dirs \
    -o ${WORK_DIR}/images/debian-bookworm-arm64/initrd.gz \
    https://deb.debian.org/debian/dists/bookworm/main/installer-arm64/current/images/netboot/debian-installer/arm64/initrd.gz
fi

# debian 13

# amd64
if [[ ! -f "${WORK_DIR}/images/debian-trixie-amd64/linux" ]]; then
    curl -s \
    --create-dirs \
    -o ${WORK_DIR}/images/debian-trixie-amd64/linux \
    https://deb.debian.org/debian/dists/trixie/main/installer-amd64/current/images/netboot/debian-installer/amd64/linux
fi

if [[ ! -f "${WORK_DIR}/images/debian-trixie-amd64/initrd.gz" ]]; then
    curl -s \
    --create-dirs \
    -o ${WORK_DIR}/images/debian-trixie-amd64/initrd.gz \
    https://deb.debian.org/debian/dists/trixie/main/installer-amd64/current/images/netboot/debian-installer/amd64/initrd.gz
fi

# arm64
if [[ ! -f "${WORK_DIR}/images/debian-trixie-arm64/linux" ]]; then
    curl -s \
    --create-dirs \
    -o ${WORK_DIR}/images/debian-trixie-arm64/linux \
    https://deb.debian.org/debian/dists/trixie/main/installer-arm64/current/images/netboot/debian-installer/arm64/linux
fi

if [[ ! -f "${WORK_DIR}/images/debian-trixie-arm64/initrd.gz" ]]; then
    curl -s \
    --create-dirs \
    -o ${WORK_DIR}/images/debian-trixie-arm64/initrd.gz \
    https://deb.debian.org/debian/dists/trixie/main/installer-arm64/current/images/netboot/debian-installer/arm64/initrd.gz
fi
