add your public key to digital ocean
get the fingerprint from there and assign that to the droplets
try to connect the droplets via that ssh key (by using an external ssh package)
try to run commands on droplet(s)

my Windows machine's ssh fingerprint : a1:15:6c:ef:9a:c3:49:80:e0:da:d5:2a:7f:41:34:00


for ssh connections, usernames are :
```
The default username is root on most operating systems, like Ubuntu and CentOS.
Exceptions to this include CoreOS, where you’ll log in as core,
Rancher, where you’ll log in as rancher,
and FreeBSD, where you’ll log in as freebsd.
```



You can get your ssh's fingerprint via this command
$ ssh-keygen -l -E md5 -f ~/.ssh/id_rsa.pub
