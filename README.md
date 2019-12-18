# Morty

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

用于消息提示，

```yaml
userlist:
  - name: zhendong
    chinesename: 振东
    phonenumber: 176****6226
```

**name:**





















