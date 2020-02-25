# Create users and bind cluster role tools
---
## Function Description
* The tool can be used without UI operations to directly create local users and bind cluster roles.

## Prerequisite
* Need to prepare API-Token；
* Need to prepare the cluster ID；
* Need to prepare the cluster role template ID；

## Using parameters

```shell
-t RancherServer Token. （Required）

-u RancherServer URL. (Required)

-t Whether to enable the user. True of false (default is true)

-m Whether the user has only login privileges. True of false        (default is false)

-c Whether to require the user to change the password upon first login. True of false (default is true)

-n show name. (Required)

-p Initial password. (Default rancher123)

-gr Create user default roles. (Default is user)

-us Create a username. (Required)

-ci Cluster ID. (Required)

-r Cluster role template id. (Required)

```

---
## Example
* Example to create a user and bind a cluster role

```shell
./Create_User_Darwin_amd64 -u 47.56.251.203 -t token-ddldl:hkjb2b9k55b769mvw4bssb6thhn6dxzbl5hrmv7vszrqpbzksd9ctf -n rancher_zhen -gr user -us rancher_zhen -ci c-8tt9d -r rt-ssrvn
```
