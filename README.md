# Morty

![Travis ](https://img.shields.io/travis/ranzhendong/morty/master?style=plastic)
![GitHub release](https://img.shields.io/github/v/release/ranzhendong/morty?include_prereleases&style=plastic)
![GitHub last commit ](https://img.shields.io/github/last-commit/ranzhendong/morty/master?style=plastic)
![GitHub](https://img.shields.io/github/license/ranzhendong/morty?style=plastic)

[Morty](#Morty)

- [介绍](#介绍)

- [配置](#配置)
  - [用户配置userlist](#用户配置userlist)
  - [k8s配置kubernetes](#k8s配置kubernetes)
  - [钉钉配置dingding](#钉钉配置dingding)

- [使用](#使用)
  - [InstantUpdate即时更新](#InstantUpdate即时更新)
    - [使用方式](InstantUpdate#使用方式)
    - [路由](#InstantUpdate路由)
    - [数据格式](#InstantUpdate数据格式)
    - [参数说明](#InstantUpdate参数说明)
    - [请求方式curl](#InstantUpdate请求方式curl)
    - [返回结果](#InstantUpdate返回结果)
    
  - [GrayUpdate混合阶梯灰度更新](#GrayUpdate混合阶梯灰度更新)
    - [使用方式](#GrayUpdate使用方式)
    - [路由](#GrayUpdate路由)
    - [数据格式](#GrayUpdate数据格式)
    - [参数说明](#GrayUpdate参数说明)
    - [请求方式curl](#GrayUpdate请求方式curl)
    - [返回结果](#GrayUpdate返回结果)
    
  - [RollBack回滚](#RollBack回滚)
    - [使用方式](#RollBack使用方式)
    - [路由](#RollBack路由)
    - [数据格式](#RollBack数据格式)
    - [参数说明](#RollBack参数说明)
    - [请求方式curl](#RollBack请求方式curl)
    - [返回结果](#RollBack返回结果)

</br></br>

## 介绍

&emsp;&emsp;Morty是以golang语言为基础而开发apiserver，借助deployment特性对kubernetes中deployement管理项目进行更新以及回滚。包含以下三种功能：

- 即时更新
- 混合阶梯灰度更新
- 回滚



</br></br>

## 配置

&emsp;&emsp;通过config.yaml对配置进行管理，内部采用[viper](https://github.com/spf13/viper)实现。

</br>

### k8s配置kubernetes用户配置userlist

```yaml
userlist:
  - name: zhendong
    chinesename: 振东
    phonenumber: 176****6226
```

- **name:** 用户简称，用于数据结构当中。

- **chinesename:** 用户中文名称，用于后面消息发送（保留字段，程序当中还暂时没有用到）。

- **phonenumber:** 用户对应的电话号码，用于钉钉@。

</br>

### k8s配置kubernetes

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



</br>

### 钉钉配置dingding

&emsp;&emsp;需要将钉钉机器人地址填写到**robotsurl**，注alertcontent意钉钉机器人有关键字过滤，防止恶意刷取消息，因此需要在源码部分[alertcontent.go](https://github.com/ranzhendong/morty/blob/master/src/public/alert/alertcontent.go)文件当中进行修改**keywords**字段。保证消息可以正常发送。



</br></br>

## 使用

&emsp;&emsp;内部是http server，因此需要通过发送请求的方式对服务进行访问以及操作。

&emsp;&emsp;下面的三个功能都是围绕以deployment为基础进行维护的pod的更新，对于其他方式启动的pod暂时不支持更新。

&emsp;&emsp;这个软件我建议结合awx（ansible web管理项目）进行使用。测试的话可以使用curl或者postman；生产环境，不建议直接通过api对服务进行请求。



</br>

### InstantUpdate即时更新

&emsp;&emsp;即时发布，顾名思义就是新版本镜像会立即生效，更新完成的时间依据replicas、minReadySeconds和pod内部健康检查参数设置，但是总的来说deployment会接管一切，更新完成。

&emsp;&emsp;等同于下面kubectl命令。

```shell
kubectl set image deployment/nginx nginx=nginx:1.9.1
```



</br>

#### InstantUpdate使用方式

&emsp;&emsp;演示我使用postman对接口进行请求。



</br>

#### InstantUpdate路由

**&emsp;&emsp;/deployupdate**



</br>

#### InstantUpdate数据格式

&emsp;&emsp;数据为json

```json
{
	"name": "InstantDeployment",
	"deployment": "nginx-deployment",
	"namespace": "default",
	"image": "nginx:1.7.9",
	"javaProject": "vmims",
	"version": "1.0.1",
	"minReadySeconds": 12,
	"replicas": 10,
	"sendFormat": "texts",
	"rollingUpdate": {
		"maxUnavailable": "40%",
		"maxSurge": "40%"
	},
	"info": {
		"requestMan": "zhendong",
		"updateSummary": "update for myself",
		"phoneNumber": "17600376226"
	}
}
```



</br>

#### InstantUpdate参数说明

**不可更改参数**

&emsp;&emsp;不能更改，否则导致程序报错。

- **name:**  url唯一标识符。



**必须更改参数**

&emsp;&emsp;必须依据自身需求进行更改的参数，否则也可能报错。

- **deployment:** 指定更新deployment名称。
- **image:** 指定更新镜像，必须保证每个节点都有这个镜像。
- **javaProject:** 指定镜像内部工程名称，用于钉钉消息提示。
- **version:** 指定镜像内部工程版本，用于钉钉消息提示。
- **sendFormat:** 指定钉钉发送消息格式。可选text（普通格式发送）；markdown格式发送
- **info:** 指定此次执行更新操作人的信息。
  - **updateSummary:** 此次更新信息。
  - **requestMan&phoneNumber: **保证和配置文件**config.yaml**一致，因为内部会进行信息比对，手机号用于钉钉@，需要保证手机号和钉钉绑定一致。



**可选更改参数**

&emsp;&emsp;可以选择不上传这些参数，保持已有默认值，不存在的话，会默认创建。

- **minReadySeconds:** deployment认为pod准备好对外接受请求时间。
- **replicas:** pod个数（副本数）
- **rollingUpdate:** 滚动更新两个参数设置，这里支持个数和百分比。
  - **maxUnavailable:** 滚动更新时最大不可用数量。
  - **maxSurge:** 滚动更新时最大可用数量。



</br>

#### InstantUpdate请求方式curl

&emsp;&emsp;可以构造curl进行请求

```shell
APISERVER=http://127.0.0.1:8080

curl $APISERVER/deployupdate -X POST -H "Content-Type:application/json" --data '{
	"name": "InstantDeployment",
	"deployment": "nginx-deployment",
	"namespace": "default",
	"image": "nginx:1.7.9",
	"javaProject": "vmims",
	"version": "1.0.1",
	"minReadySeconds": 12,
	"replicas": 10,
	"sendFormat": "texts",
	"rollingUpdate": {
		"maxUnavailable": "40%",
		"maxSurge": "40%"
	},
	"info": {
		"requestMan": "zhendong",
		"updateSummary": "update for myself",
		"phoneNumber": "17600376226"
	}
}'
```



</br>

#### InstantUpdate返回结果

当返回结果为下面信息，说明成功。

如果失败会返回响应的错误原因。

```text
[Main.DpUpdate] Deployment Image Update Complete!
```



</br></br>

### GrayUpdate混合阶梯灰度更新

</br>



#### GrayUpdate使用方式

&emsp;&emsp;演示我使用postman对接口进行请求。



</br>

#### GrayUpdate路由

**&emsp;&emsp;/graydeployupdate**



</br>

#### GrayUpdate数据格式

&emsp;&emsp;数据为json

```json
{
	"name": "GrayDeployment",
	"deployment": "nginx-deployment",
	"namespace": "default",
	"image": "nginx:1.9.1",
	"javaProject": "vmims",
	"version": "1.0.1",
	"minReadySeconds": 5,
	"replicas": 4,
	"sendFormat": "texts",
	"gray": {
		"tieredRate": "15%",
		"durationOfStay": "3min",
		"aVersionStepWiseUp": "55s",
		"bVersionStepWiseUp": "45s",
		"bVersionStepWiseDown": "37s"
	},
	"rollingUpdate": {
		"maxUnavailable": "40%",
		"maxSurge": "40%"
	},
	"info": {
		"requestMan": "zhendong",
		"updateSummary": "update for myself",
		"phoneNumber": "17600376226"
	}
}
```

</br>



#### GrayUpdate参数说明

**不可更改参数**

&emsp;&emsp;不能更改，否则导致程序报错。

- **name:**  url唯一标识符。



**必须更改参数**

&emsp;&emsp;必须依据自身需求进行更改的参数，否则也可能报错。

- **deployment:** 指定更新deployment名称。
- **image:** 指定更新镜像，必须保证每个节点都有这个镜像。
- **javaProject:** 指定镜像内部工程名称，用于钉钉消息提示。
- **version:** 指定镜像内部工程版本，用于钉钉消息提示。
- **sendFormat:** 指定钉钉发送消息格式。可选text（普通格式发送）；markdown格式发送
- **gray:** 指定灰度发布参数
  - **tieredRate:** 灰度发布比例。例如replicas=10，tieredRate=20%那么发布比例将按照2，阶梯增加，逐渐增加pod个数。
  - **durationOfStay:** 灰度发布持续时间。当新版本启动到指定数目，**新老版本共存时间**。
  - **aVersionStepWiseUp:** A版本阶梯增加持续时间。**老版本每增加tieredRate个pod**，持续aVersionStepWiseUp时间，然后再次增加tieredRate个pod。
  - **bVersionStepWiseUp:** B版本阶梯增加持续时间。**新版本每增加tieredRate个pod**，持续bVersionStepWiseUp时间，然后再次增加tieredRate个pod。
  - **bVersionStepWiseDown:** B版本阶梯递减持续时间。**老版本在镜像更新之后，就将最先启动的新版本每递减tieredRate个pod**，持续bVersionStepWiseDown时间，然后再次递减tieredRate个pod。
- **info:** 指定此次执行更新操作人的信息。
  - **updateSummary:** 此次更新信息。
  - **requestMan&phoneNumber: **保证和配置文件**config.yaml**一致，因为内部会进行信息比对，手机号用于钉钉@，需要保证手机号和钉钉绑定一致。



**可选更改参数**

&emsp;&emsp;可以选择不上传这些参数，保持已有默认值，不存在的话，会默认创建。

- **minReadySeconds:** deployment认为pod准备好对外接受请求时间。
- **replicas:** pod个数（副本数）
- **rollingUpdate:** 滚动更新两个参数设置，这里支持个数和百分比。
  - **maxUnavailable:** 滚动更新时最大不可用数量。
  - **maxSurge:** 滚动更新时最大可用数量。



</br>

#### GrayUpdate请求方式curl

&emsp;&emsp;可以构造curl进行请求

```shell
APISERVER=http://127.0.0.1:8080

curl $APISERVER/deployupdate -X POST -H "Content-Type:application/json" --data '{
	"name": "GrayDeployment",
	"deployment": "nginx-deployment",
	"namespace": "default",
	"image": "nginx:1.9.1",
	"javaProject": "vmims",
	"version": "1.0.1",
	"minReadySeconds": 5,
	"replicas": 4,
	"sendFormat": "texts",
	"gray": {
		"tieredRate": "15%",
		"durationOfStay": "3min",
		"aVersionStepWiseUp": "55s",
		"bVersionStepWiseUp": "45s",
		"bVersionStepWiseDown": "37s"
	},
	"rollingUpdate": {
		"maxUnavailable": "40%",
		"maxSurge": "40%"
	},
	"info": {
		"requestMan": "zhendong",
		"updateSummary": "update for myself",
		"phoneNumber": "17600376226"
	}
}'
```



</br>

#### GrayUpdate返回结果

当返回结果为下面信息，说明成功。

如果失败会返回响应的错误原因。

```text
[Main.GrayDeploymentUpdate] Deployment Image Update Complete!
```

</br>



### RollBack回滚



</br>

#### RollBack使用方式

&emsp;&emsp;演示我使用postman对接口进行请求。



</br>

#### RollBack路由

**&emsp;&emsp;/rollback**



</br>

#### RollBack数据格式

&emsp;&emsp;数据为json

```json
{
	"name": "RollBack",
	"deployment": "nginx-deployment",
	"namespace": "default",
	"image": "nginx:1.7.9",
	"javaProject": "vmims",
	"version": "1.0.1",
	"sendFormat": "texts",
	"info": {
		"requestMan": "zhendong",
		"updateSummary": "update for myself",
		"phoneNumber": "17600376226"
	}
}
```

</br>



#### RollBack参数说明

**不可更改参数**

&emsp;&emsp;不能更改，否则导致程序报错。

- **name:**  url唯一标识符。



**必须更改参数**

&emsp;&emsp;必须依据自身需求进行更改的参数，否则也可能报错。

- **deployment:** 指定更新deployment名称。
- **image:** 指定更新镜像，必须保证每个节点都有这个镜像。
- **javaProject:** 指定镜像内部工程名称，用于钉钉消息提示。
- **version:** 指定镜像内部工程版本，用于钉钉消息提示。
- **sendFormat:** 指定钉钉发送消息格式。可选text（普通格式发送）；markdown格式发送
- **info:** 指定此次执行更新操作人的信息。
  - **updateSummary:** 此次更新信息。
  - **requestMan&phoneNumber: **保证和配置文件**config.yaml**一致，因为内部会进行信息比对，手机号用于钉钉@，需要保证手机号和钉钉绑定一致。



</br>

#### RollBack请求方式curl

&emsp;&emsp;可以构造curl进行请求

```shell
APISERVER=http://127.0.0.1:8080

curl $APISERVER/deployupdate -X POST -H "Content-Type:application/json" --data '{
	"name": "RollBack",
	"deployment": "nginx-deployment",
	"namespace": "default",
	"image": "nginx:1.7.9",
	"javaProject": "vmims",
	"version": "1.0.1",
	"sendFormat": "texts",
	"info": {
		"requestMan": "zhendong",
		"updateSummary": "update for myself",
		"phoneNumber": "17600376226"
	}
}'
```



</br>

#### RollBack返回结果

当返回结果为下面信息，说明成功。

如果失败会返回响应的错误原因。

```text
[Main.RollBack] Deployment Rollback Successful!
```



</br></br>

# Copyright & License

BSD 2-Clause License

Copyright (c) 2019, Zhendong
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

- Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

- Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.