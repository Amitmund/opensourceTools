## Purpose
To have a map between a pid to a docker container if any.

---

## To build:
go build main.go

---
## Help:
```
Usage:
        This program expect one and only one numerical pid number as an command-line argument.
        So that, it will try to map if the given pid# belongs to any container...

Example:
        pidToContainerCheck pid_number
        pidToContainerCheck 1234


QuickNote:
        This tool best goes with the following command output...
        ps -eo pcpu,pid,user,size,priority,args | sort -k1 -r -n | head -10
```


---

## sample output:

```
./pidToContainerCheck 210682
Found the given PID belongs to: container_id:  450b4e7bce15

Applicaiton Image Name: 'nginx'
Image Maintainer: 'NGINX Docker Maintainers <docker-maint@nginx.com>'
Container Started At: '2022-01-11T15:00:32.720704386Z'
Container Shared Memory in Bytes: '67108864'

You might like to run  docker container top 450b4e7bce15  to see further details.
You might like to run  docker ps | grep 450b4e7bce15  to see further details.
You might like to run  docker stats 450b4e7bce15  to see further details.
```

---
