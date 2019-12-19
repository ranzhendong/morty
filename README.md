# Morty

![Travis (README.assets/master.svg) branch](https://img.shields.io/travis/ranzhendong/morty/master?style=plastic)
![GitHub release (README.assets/lecter-1576747587037.svg)](https://img.shields.io/github/v/release/ranzhendong/morty?include_prereleases&style=plastic)
![GitHub last commit (README.assets/master-1576747587023.svg)](https://img.shields.io/github/last-commit/ranzhendong/morty/master?style=plastic)
![](README.assets/morty.svg)

[Morty](#Morty)

- [介绍](#介绍)
- [配置](#配置)
  - [用户配置userlist](#用户配置userlist)
  - 

## 介绍

&emsp;&emsp;Morty是以golang语言为基础而开发apiserver，借助deployment特性对kubernetes中deployement管理项目进行更新以及回滚。包含以下三种功能：

- 即时更新
- 混合阶梯灰度更新
- 回滚

## 用户配置userlist配置

&emsp;&emsp;通过config.yaml对配置进行管理，内部采用[viper](https://github.com/spf13/viper)实现。



#### 用户配置userlist

```yaml
userlist:
  - name: zhendong
    chinesename: 振东
    phonenumber: 176****6226
```

- **name:** 用户简称，用于数据结构当中。

- **chinesename:** 用户中文名称，用于后面消息发送（保留字段，程序当中还暂时没有用到）。

- **phonenumber:** 用户对应的电话号码，用于钉钉@。



#### k8s配置kubernetes

```yaml
kubernetes:
  host: https://172.16.0.60:6443
  tokenfile:
  deploymentapi: /apis/extensions/v1beta1
```

- **host:** k8smaster节点地址和端口，如果在集群内部允许，可以将地址改成内部域名，访问速度以及安全性提高。

因为在win上开发的项目，所以配置了从集群外部进行访问。

- **tokenfile:** 可以将访问token写在这里，也可以选择使用外部token文件，但是需要保证外部token文件是json格式。

目前还未解决直接读取挂载serviceaccount的token文件。

- **deploymentapi:** 因为k8s版本不同，导致api版本也不尽相同，因此把这部分作为配置，尽量保证多版本k8s可以运行。



#### 钉钉配置dingding

需要将钉钉机器人地址填写到**robotsurl**，注意需要将关键字进行过滤



