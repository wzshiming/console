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

### SSH  

> http://localhost:8888/?name=ssh&host=ssh://{username}:{password}@{domain}:{port}  

### Docker  

> http://localhost:8888/?name=docker&host=tcp://localhost:2375&cid={Container}&cmd={cmd}  
or  
> http://localhost:8888/?name=docker&host=/var/run/docker.sock&cid={Container}&cmd={cmd}  

## MIT License

Copyright © 2017-2018 wzshiming<[https://github.com/wzshiming](https://github.com/wzshiming)>.

MIT is open-sourced software licensed under the [MIT License](https://opensource.org/licenses/MIT).