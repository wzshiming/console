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

``` bash
$(go env GOBIN)/web_console
```
or 
``` bash
docker run -it --rm -p 8888:8888 wzshiming/web_console
```

### Shell  

> http://localhost:8888/?name=shell&cmd={cmd}  

### SSH  

> http://localhost:8888/?name=ssh&host=ssh://{username}:{password}@{domain}:{port}  

### Docker  

> http://localhost:8888/?name=docker&host=tcp://localhost:2375&cid={Container}&cmd={cmd}  
or  
> http://localhost:8888/?name=docker&host=unix:///var/run/docker.sock&cid={Container}&cmd={cmd}  

## License

Pouch is licensed under the MIT License. See [LICENSE](https://github.com/wzshiming/console/blob/master/LICENSE) for the full license text.