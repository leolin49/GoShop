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

### time
1. 新增 time-service 用于定时任务，如限时秒杀等

### stock
1. 新增 基于redis缓存预热 + Lua脚本 的秒杀实现
2. TODO: 秒杀结束后，缓存回填到数据库

## 2025.03.02

### login
1. 分布式session设计？jwt是否应该放在session里？(不应该，因为jwt是无状态的设计，和session相违背)
2. TODO: 最终方案：
    2.1 用户登录时，返回短token和sessionId，并将session信息存入redis
    2.2 对于每个请求，gin鉴权中间件先判断短token是否过期
        如果过期了，则根据sessionId从redis中拿到长token，验证长token是否过期
            2.2.1 如果没过期，则重新生成双token，更新session信息到redis中，完成续签
            2.2.2 如果过期了，删除session，返回错误码表示需要重新登录

### stock
1. TODO: 秒杀和普通购买是不是走同一套逻辑？库存是否需要同步？
2. 如何判断当前是否处于秒杀活动期间？

## 2025.03.02

### redis
1. 利用分布式锁实现部分接口的幂等性（支付）

> 已接入redis缓存的服务: cart checkout

## 2025.03.27

### login
1. 注册成功后发送邮件（TODO）

### share 
1. 长链接转换成短链接 (TODO)
    1. 短链接生成
    2. 使用布隆过滤器快速判断短链接是否存在，避免频繁访问数据库
    方案1: 使用hash函数生成，如果产生哈希冲突则在原链接加上特定标记字符串，再进行hash
        e.g. goshop.com/aaaa-bbbb-cccc-dddd -- hash -> goshop.com/A8E3B4
         1) 如果发现A8E3B4已存在，则判断对应的原链接是否一样，如果一样则直接返回
         2) 如果不一样，则表明发生了hash冲突，在原链接后添加 "!@#$%^&"
         3) aaaa-bbbb-cccc-dddd!@#$%^& -- hash -> goshop.com/C762BA
         4) 重复步骤1，直到没有产生冲突为止
         5) 查询短链接时，拿到长链接后，需要先去掉(如果有的话)后面的"!@#$%^&"
    方案2: 每次预生成一个短链接池，然后分配（一个长链接可能对应多个短链接）
         1) 每天生成10w个UUID，分配完成后在数据库中维护其对应关系 


