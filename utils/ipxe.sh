#!/usr/bin/env bash

WORK_DIR=$(pwd)

# apt install -y make gcc-aarch64-linux-gnu liblzma-dev

git submodule init
git submodule update
git submodule

cd ${WORK_DIR}/ipxe/src

# make CROSS=aarch64-linux-gnu- bin-arm64-efi/ipxe.efi EMBED=ipxe.script

# arm64 uefi
if [[ ! -f "${WORK_DIR}/tftpboot/ipxe-arm64.efi" ]]; then
    make CROSS=aarch64-linux-gnu- bin-arm64-efi/ipxe.efi
    cp bin-arm64-efi/ipxe.efi ${WORK_DIR}/tftpboot/ipxe-arm64.efi
fi

# x86_64 uefi
if [[ ! -f "${WORK_DIR}/tftpboot/ipxe-x86_64.efi" ]]; then
    make bin-x86_64-efi/ipxe.efi
    cp bin-x86_64-efi/ipxe.efi ${WORK_DIR}/tftpboot/ipxe-x86_64.efi
fi

# x86_64 bios
if [[ ! -f "${WORK_DIR}/tftpboot/ipxe-x86_64.pxe" ]]; then
    make bin-x86_64-pcbios/ipxe.pxe
    cp bin-x86_64-pcbios/ipxe.pxe ${WORK_DIR}/tftpboot/ipxe-x86_64.pxe
fi

cd ${WORK_DIR}
