# web console  

## Download & Install 

``` bash
go get -u -v github.com/wzshiming/console/cmd/web_console
```
or  
``` 
docker pull wzshiming/web_console 
```

## Start 

```
docker run -it --rm -p 8888:8888 wzshiming/web_console
```

### SSH  

> http://localhost:8888/?name=ssh&host=ssh://{username}:{password}@{domain}:{port}  

### Docker  

> http://localhost:8888/?name=docker&host=tcp://localhost:2375&cid={Container}&cmd={cmd}  
or  
> http://localhost:8888/?name=docker&host=/var/run/docker.sock&cid={Container}&cmd={cmd}  

