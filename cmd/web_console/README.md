# web console



## ssh

> http://localhost:8888/?name=ssh&host=ssh://{{username}}:{{password}}@{{domain}}:{{port}}

## docker

> http://localhost:8888/?name=docker&host=tcp://localhost:2375&cid={{Container}}&cmd={{cmd}}
or
> http://localhost:8888/?name=docker&host=/var/run/docker.sock&cid={{Container}}&cmd={{cmd}}

## shell

> http://localhost:8888/?name=shell&cmd=sh