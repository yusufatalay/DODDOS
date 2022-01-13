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


# TODOS
- [x] Be able to create and boot machines through digital ocean API-------------------------------------------------------------> DONE
- [x] Creating can be done with OR without tags---------------------------------------------------------------------------------> DONE
- [x] Deleting can be done with ID or TAG---------------------------------------------------------------------------------------> DONE
- [x] Be able to delete machine-------------------------------------------------------------------------------------------------> DONE
- [x] Show machine info (IP, name, id etc.)-------------------------------------------------------------------------------------> DONE
- [x] Be able to list all droplets or droplets that has same tag----------------------------------------------------------------> DONE
- [x] Be able to create multiple droplets at once (bind them with a tag)--------------------------------------------------------> DONE
- [x] Be able to define a SSH key to created machine----------------------------------------------------------------------------> DONE
- [x] Be able to restart machine------------------------------------------------------------------------------------------------> DONE
- [x] Be able to operate created machines via ssh (make them send requests to the target etc.)----------------------------------> DONE
- [x] Add this project to github with a proper documentation and add "regrews" as a contributer IMPORTANT-----------------------> DONE
- [ ] Add a killswitch (commad is "killall -u root")  IMPORTANT
- [ ] Be able to restart machines with tag name		 IMPORTANT
- [ ] Be able define ssh key inside of the code		 IMPORTANT
- [X] Be able to run scripts on created machines		 IMPORTANT-----------------------------------------------------------------> DONE (just upload the code to the github and curl from it)
- [ ] Be able to see only one machine's output with flag  'logging mechanism'
- [ ] Be able to see all machine's outputs with another flag 'logging mechanism'
- [ ] Try to make a WEB interface with react
- [ ] Get the bandwith metrics of the all the machines merged together pipe the output to a csv file in the format of (time:bandwith)
- [ ] Get the output form above and make a graph in matplotlib
