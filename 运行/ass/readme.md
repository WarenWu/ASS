# asgateway打包
```
.
├── bin					// 执行程序与动态库目录，动态库暂时先放这里
├── conf				// 配置文件路径
├── log					// 日志路径(运行后会自动创建)
├── install.bat
├── start.bat
├── status.bat
├── stop.bat
├── uninstall.bat
├── asgateway_svr.exe		// 服务管理程序(winsw)
└── asgateway_svr.xml		// 服务管理配置文件
```
**在xml可配置服务PATH环境变量，经过测试在windows 11上有效，windows server 2008无效，所以动态库暂时放bin目录下**

### winsw
* download: https://github.com/winsw/winsw  
* version: 2.10.3