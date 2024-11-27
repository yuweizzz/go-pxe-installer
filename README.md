## Get started

通过网络启动，在虚拟机环境下快速安装 Debian 12 。

## Build and Run

执行构建。

```bash
# required go version go1.21.0 and make

# no embed images
make ipxe
make build

# embed images
make ipxe
make images
make buildi
```

修改配置后可以直接运行。

```yaml
# tftp example: no embed images
# 修改成实际机器的网卡
iface: enp0s3
# 设置实际的 IP 地址
ipaddr: 10.0.2.5
logger:
  level: debug
  file: /dev/stdout
dhcp:
  port: 67
tftp:
  port: 69
  # 使用外部的 tftp 目录，如果在外部目录无法打开文件，会尝试在内嵌目录搜索
  external: help
pxe:
  # default menu target, could be 'config', 'shell', 'reboot', 'exit' or entries label
  default: shell
  # timeout(Unit: ms), '0' means not auto choose the default option
  timeout: 0
  entries:
    - display: Debian 12 bookworm amd64
      label: x86_64
      kernel: tftp://10.0.2.5/images/debian-bookworm-amd64/linux
      initrd: tftp://10.0.2.5/images/debian-bookworm-amd64/initrd.gz
      # if use tftp in preseed, the tftp server should follow the value of ipaddr
      # or use http like this: preseed/url=http://somewhere/preseed.txt
      append: initrd=initrd.gz vga=normal fb=false auto=true priority=critical preseed/url=tftp://10.0.2.5/debian12-preseed.txt
    - display: Debian 12 bookworm arm64
      label: arm64 
      kernel: tftp://10.0.2.5/images/debian-bookworm-arm64/linux
      initrd: tftp://10.0.2.5/images/debian-bookworm-arm64/initrd.gz
      append: initrd=initrd.gz vga=normal fb=false auto=true priority=critical preseed/url=tftp://10.0.2.5/debian12-preseed.txt
```

### Orcale VM VirtualBox

> [!IMPORTANT]
> 1. 网络启动后会停留在镜像选择步骤，选择对应镜像后会在安装过程对硬盘进行格式化，对于已有系统的虚拟电脑应该慎重操作。
> 2. 当前的 Debian 12 安装镜像默认设置普通用户 debian 和 root 用户两个，密码和各自的用户名相同，系统安装完成后请务必修改。

新建虚拟电脑 -> 设置 -> 系统 -> 启动顺序 -> 设置为网络优先。

根据当前服务的运行位置决定网络连接方式：

- 如果已经有多台虚拟电脑运行在同个 NAT 网络下，可以直接将网络设置成 NAT 网络，并在其中任意一台虚拟电脑运行本服务。
- 如果虚拟电脑是独立且全新的，那么直接选择桥接网络，后续将本服务运行在宿主机即可。

安装完成后需要重新调整启动顺序为硬盘优先。

## Todo

- [x] ~~当前的 pxelinux.cfg/default 文件需要渲染 tftp 地址，否则应该手动修改后重新编译。~~
- [x] ~~允许镜像从远程拉取~~。
- [x] ~~指定额外的本地 tftp 目录~~。
- [x] ~~解决 ESXi 不兼容的问题~~。
- [x] ~~解决 UEFI 环境下远程下载不稳定的问题~~。
- [ ] 移除 syslinux ，通过 ipxe 实现控制。

## Where file from

### PXE boot (deprecated)

以下所有文件都基于已经安装完成的 Debian 12 系统。

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

通过 syslinux/pxelinux 支持的 bootfile 已经弃用，使用 iPXE 代替。

### Images

debian-bookworm-amd64:

``` shell
# linux
wget https://deb.debian.org/debian/dists/bookworm/main/installer-amd64/current/images/netboot/debian-installer/amd64/linux
# initrd.gz
wget https://deb.debian.org/debian/dists/bookworm/main/installer-amd64/current/images/netboot/debian-installer/amd64/initrd.gz
```

debian-bookworm-amd64:

``` shell
# linux
wget https://deb.debian.org/debian/dists/bookworm/main/installer-arm64/current/images/netboot/debian-installer/arm64/linux
# initrd.gz
wget https://deb.debian.org/debian/dists/bookworm/main/installer-arm64/current/images/netboot/debian-installer/arm64/initrd.gz
```

### Others

* `example-preseed.txt`: download from [d-i.debian.org](https://d-i.debian.org/manual/example-preseed.txt), `debian12-preseed.txt` is modified from it.
* `example-ipxe.script`: iPXE script example.
