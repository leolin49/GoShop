# PROCESS LOG

## 2025.02.25

### RabbitMQ
1. 接入到了gateway和checkout服务之间，实现了MQ限流和削峰；
2. TODO: 在checkout和pay，order等服务之间接入MQ，以实现异步通信和服务间解耦合；

### Consul 
1. gateway实现了从consul配置中心读取配置； 2. TODO: 需要实现consul长轮询读取配置，以实现热更新；
3. TODO: 重构pkg/service/consul.go，将consul相关封装；(DONE)
4. TODO: 其他服务也需要改为从consul读取配置；(DONE)

### Mysql
1. TODO: Mysql连接改为由pkg/mysql/db.go统一生成，只需要传入配置；(DONE)

### Redis
1. cartserver实现了GetCart缓存（注意先更新Mysql再删除Redis缓存以保证数据一致性），在AddCart和CleanCart时需要删除缓存；
2. TODO: 其他需要做缓存的接口；

### Gateway
1. 通过启动参数区分不同节点，并在Nginx实现Gateway的负载均衡；

## 2025.02.26

### GoMock
1. 通过gomock实现redis模块的单测，需要强制使用面向interface编程；
2. TODO: 其他模块也通过mock测试？

### .env
1. TODO: 地址好像有问题

### Mysql
1. 实现了Mysql主从复制的配置和docker环境(DONE)
2. TODO: 使用gorm实现读写分离(DONE)

## 2025.02.27

### login
1. TODO: mysql读写分离 consul配置中心(DONE)

### 代码重构
2. gateway login product auth

## 2025.02.28

### login
1. 密码做md5盐值加密(DONE)

### stock
1. TODO: 新增库存服务(DONE)
2. 基于Mysql事务和乐观锁防止超卖（但是高并发下，行级锁存在性能问题）

### redis
1. TODO: 部署redis集群
2. TODO: 基于redis实现分布式锁

### 代码重构
1. all service done 


## 2025.03.01

### grpc
1. TODO: retry + 接口幂等性设计


> 已接入redis缓存的服务: cart checkout


