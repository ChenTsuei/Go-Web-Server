<p>这个项目是我的计算机网络课程的大作业，使用Go语言编写，行数在200行左右，受到了开源项目<a href="http://tinyhttpd.sourceforge.net/">Tiny HTTPd</a>的许多启发。</p>
<p>实现了简单的HTTP方法：GET / POST / HEAD，并支持CGI。</p>
<p>项目地址：<a href="https://github.com/ChenTsuei/Go-Web-Server">https://github.com/ChenTsuei/Go-Web-Server</a></p>

## 使用方法

```
git clone https://github.com/ChenTsuei/Go-Web-Server.git
cd ./Go-Web-Server
go build gws.go
./gws --port=[端口号] --path=[根路径]
```
