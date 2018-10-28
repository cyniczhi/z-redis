# z-redis

## Intro
https://gist.github.com/tevino/f526b18117d83a0fb7b4d1857b14051a
基于Golang从头编写的一个类似Redis的Demo，目前只支持string类型的value。
## Preparation
### get/set/del
基本思路是使用Go自带的数据类型Map作为底层存储的数据结构。
注意:
1. 似乎Map只支持增量扩容，不支持缩容。所以Map有可能会一直增长，内存得不到释放。因此需要考虑实现可以shrink的哈希表。
2. Map是非线程安全的，因此多用户操作同一个db有可能会出现问题。
### redis标准协议
> http://redisdoc.com/topic/protocol.html

### LRU过期
对Database新增expire字典，key和dict里面的key一一对应，value是一个链表的结点。
可以为不同的过期策略实现不同的expire字典，这里只考虑对LRU的实现。

#### expire字典:
插入：新增key的时候，加入到链表尾部。检查是否超过最大链表长度，如果超过就删除表头。
更新：每当key被使用（set/get），把key对应的node移动到链表尾部。
删除：删除key的时候，把key对应的node删除，再删除expire里面的key。
**LRUDict:**
- Dict map[string]*Node
  - 缓存key的字典
- Max  int32
  - 缓存最多多少key
- Head *Node
  - 链表头
- Tail *Node
  - 链表尾
- Len  int32
  - 链表长度


链表:
双向无环链表，包含头指针和尾指针。
**Node**：
- Prev       *Node
- Next       *Node
- Key        string
  - 对应的到LRUDict的Key值
- ExpireTime int64
  - 保留字段，以便加上对过期时间的支持。

### 持久化
####
初步了解了Redis的持久化策略，本项目打算参考RDB持久化实现一个类似的ZDB持久化：将键值对按照持久化策略保存到ZDB文件里面。
####
这里没有对SAVE等主动保存指令提供支持，仅支持默认情况下的自动保存；不支持手动加载的指令及操作，只支持初始化服务器时加载ZDB文件。
## Implementation
### Client
和服务器建立连接，对Redis指令根据redis标准协议做简单封装，发送指令并展示结果。
### Server
#### ZDB持久化
##### 保存条件
给server加上saveparams字段，每个saveparam结构包含了seconds，change两个字段。即每次定期检查是否持久化的时候，看seconds秒内是否发生了大于change数量的操作（先考虑set，get，del所有），如果是就执行持久化。
为了完成上述分析，再给server加上dirty字段和lastsave字段，分别表示change数量和上次持久化的时间。
##### ZDB文件分析
先看RDB文件结构：
- RDB文件
| REDIS | db_version | databases | EOF | checksum |

- databases部分：
| REDIS | db_version | database0 | database3 | EOF | checksum |
- 每个非空数据库：
| SELECTDB | db_number | key_value_pairs |
SELECTDB 是一个1字节的常量，表示后面是一个db_number。
- key_value_pairs:
  | TYPE | KEY | VALUE | or  | EXPIRETIME_MS | ms | TYPE | KEY | VALUE |
- Value的编码
  当前z-redis的实现仅支持string类型的value，不考虑压缩，因此没有编码，仅针对value作变长地完整保存。

经过参考和分析，得到ZDB文件结构如下：
**ZDB**
- ZDB文件
	- | ZREDIS | db_version | databases | EOF |
		- ZREDIS: 'Z', 'R', 'E', 'D', 'I', 'S'; 5byte
		- db_version: 0; 1byte
		- checksum的细节日后再说，当前版本不添加。

- databases:
	 - | database0 | database3 | ... |
	 - database: | SELECTDB | db_number | key_value_pairs |
		 - SELECTDB: 1, 1 byte
		 - db_number: number, 1byte

- key_value_pairs
	- 此时先不考虑过期时间，因此只保存kv对
	- | TYPE | KEY | VALUE |
		- TYPE: const [0, 1, ..., n]: 1byte (先只支持string类型，因此TYPE全部为0)
		- KEY: string, 变长
		- VALUE: string, 变长

- value编码
	- string: 原样存储，包含'\n'
	- 其它：以后实现

***由于要关注的细节比较多，考虑了多DB的实现，并且对Golang的IO不太熟悉，这一部分多花了一些时间。***
#### 数据结构
具体参见源码及注释。
#### 网络
目前只使用了Golang默认的net包的socket编程，没有进一步考虑更高效的IO模型。
## Feature List
### 指令
- set key val
- get key
- del key
### 其它
- LRU过期策略
- 捕捉到server退出信号时会执行持久化
- server启动时会从ZDB文件里面加载数据
	- **注意：** 此时会构建一个默认的LRU字典，由于Golang map是无序的，可能会和持久化之前的LRU字典不一致。

## TODO
以下还**没有**或者只有**部分**实现。
### 指令
- select
	- 客户端多db的切换
- save
	- 客户端主动发起持久化请求
### 其它
#### 特性
- Redis统一协议。
- 更灵活的持久化策略。
- 更高效的IO模型 （网络，持久化...）
- 更少的内存开销 （支持shrink的dict）
- 更灵活的配置项启动server，启动参数/环境变量/配置文件
	- ZDB file path
	- 过期策略相关参数
	- ...
#### 测试
- 测试用例
- benchmark及性能分析
## Install
### Client
由于比较粗糙，目前客户端只起到一个发送指令和展示响应结果的作用，可以直接用telnet替代。
```
// cd to client's main package
cd ./client
go build -o client.out main.go

// by default, client will connect to 127.0.0.1:9999
./client.out
```
**OR**
```
telnet <dst addr> <dst port>
```
### Server
默认ZDB文件为default.zdb, 监听在0.0.0.0:9999，最多缓存5个k, v对，超出的按照LRU的策略过期。
```
// cd to server's main package, which is the root path of the repo by default.

go build -o server.out main.go
./server.out
```
##  Reference

> Golang的map实现

>《Redis设计与实现》chap3 chap9 chap10

> https://segmentfault.com/a/1190000015224870