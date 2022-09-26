# mindblogEnd

##格式
 api  抛出api接口

 config 各数据库以及配置文件

 core和lib 常用库的封装

 doc swagger 封装

 model 全部模型，数据库由该模型生成

 middleware 中间键

 log 存储 zap的目录 （后期要丢到数据库上）

 resource 初始化的数据都在里面

 router 每一个模块的对外抛出
 
 server 服务具体实现

## run 起来
```
 go mod tidy
 //因为为了封装  jwt.NewNumericDate 出的幺蛾子，并且又想放到 swagger,已测试到parseDepth的最小层
 swag init --parseDependency --parseInternal --parseDepth 5  -g main.go  &&  go run main.go
``` 
![image](https://user-images.githubusercontent.com/27894531/192190385-a1a7062d-e987-472b-8a44-6f117c7e8f41.png)

就可以run 起来了
