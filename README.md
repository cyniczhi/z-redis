# z-redis
@(自己造轮子)[lab]

[toc]

## Intro
https://gist.github.com/tevino/f526b18117d83a0fb7b4d1857b14051a
基于Golang从头编写的一个类似Redis的Demo，目前只支持string类型的value。
## Prep
### get/set/del
基本思路是使用Go自带的数据类型Map作为底层存储的数据结构。
注意:
1. 似乎Map只支持增量扩容，不支持缩容。所以Map有可能会一直增长，内存得不到释放。参考项目并没有做shrink。
2. Map是非线程安全的，参考项目并没有做保护。
3. 似乎zobject存储的都是512大小存储区
### redis标准协议
先作一个简单的了解，对标准协议的支持后面实现。
### LRU过期
数组+链表？
### 持久化
初步了解了Redis的持久化策略，本项目打算参考RDB持久化实现一个类似的ZDB持久化：将键值对按照持久化策略保存到ZDB文件里面。这里不对SAVE等主动保存指令提供支持，仅支持默认情况下的自动保存；不支持手动加载的指令及操作，只支持初始化服务器时加载ZDB文件。
## Implementation
### Client
和服务器建立连接，对Redis指令根据redis标准协议做简单封装，发送指令并展示结果。
### Server
#### ZDB持久化
##### 保存条件
给server加上saveparams字段，每个saveparam结构包含了seconds，change两个字段。即每次定期检查是否持久化的时候，看seconds秒内是否发生了大于change数量的操作（先考虑set，get，del所有），如果是：执行持久化。
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

#### 数据结构
#### 网络模型
## TODO
- 网络
	- IO模型
	- Buffer
- Test
	- test cases
	- benchmark
## Doc
### 安装服务器
### Demo

> hashMap of go src
> https://segmentfault.com/a/1190000015266971
> redis设计与实现