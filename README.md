apt install pxelinux
apt install syslinux

apt install syslinux-efi
cp /usr/lib/SYSLINUX.EFI/efi64/syslinux.efi .
cp /lib/syslinux/modules/efi64/ldlinux.e64 .


cp  /usr/lib/PXELINUX/pxelinux.0 .
cp /usr/lib/syslinux/modules/bios/ldlinux.c32 .