## Where file from

### PXE boot

All files based on Debian 12.

BIOS:

* `pxelinux.0`
* `ldlinux.c32`

``` shell
apt install syslinux
apt install pxelinux
cp /usr/lib/syslinux/modules/bios/ldlinux.c32 tftpboot/ldlinux.c32
cp /usr/lib/PXELINUX/pxelinux.0 tftpboot/pxelinux.0
```

UEFI:

* `syslinux.efi`
* `ldlinux.e64`

``` shell
apt install syslinux
apt install syslinux-efi
cp /lib/syslinux/modules/efi64/ldlinux.e64 tftpboot/ldlinux.e64
cp /usr/lib/SYSLINUX.EFI/efi64/syslinux.efi tftpboot/syslinux.efi
```

### Images

debian-bookworm-amd64:

* `linux`
* `initrd.gz`

``` shell
wget http://http.us.debian.org/debian/dists/bookworm/main/installer-amd64/current/images/netboot/netboot.tar.gz
tar -xvf netboot.tar.gz -C netboot
cp netboot/debian-installer/amd64/linux tftpboot/images/debian-bookworm-amd64/linux
cp netboot/debian-installer/amd64/initrd.gz tftpboot/images/debian-bookworm-amd64/initrd.gz
```

`preseed.cfg` is modified from `example-preseed.txt`.

### Others

* `example-preseed.txt`: download from [d-i.debian.org](https://d-i.debian.org/manual/example-preseed.txt)
