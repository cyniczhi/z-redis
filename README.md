# z-redis

[toc]

## Intro
https://gist.github.com/tevino/f526b18117d83a0fb7b4d1857b14051a
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
## Implementation
### Client
### Server
#### 数据结构
#### 网络模型
## TODO
- 网络
	- IO模型
	- Buffer

