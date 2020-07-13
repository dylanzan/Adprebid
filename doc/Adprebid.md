# Adprebid

## prebid现况

&emsp;&emsp;各广告媒体请求量多，Tomcat处理能力达到一定瓶颈,经常出现宕机现象。
    
## 需求分析

&emsp;&emsp;为解决当前的tomcat瓶颈问题，将各媒体请求流量按照dealID,进行流量分发。

## 功能要求

## 业务流程图

![avator](https://raw.githubusercontent.com/mindzilla/Adprebid/master/doc/Adprebid_fc.png?token=AFUG7MRYQNAG6IZDHPHJV727A7IC6 "fc")

## 业务流程

&emsp;&emsp;当媒体流量进入adprebid后，首先经过router,router根据请求url path路径，将流量路由到指定媒体controller。<br>
流量进入指定mediaCtl后，控制层会将request body,进行解析，并按照新的格式进拼装，body拼装完成后，将请求转发到service层，<br>
在controller层,controller会根据，body中的deal ID，将流量转发到配置文件中，指定的adkit服务器中。<br>
流量在adkit中，做逻辑处理，处理完成后，将response，传递给controller层，并将response body二次处理，并输出。<br>
至此，媒体请求流量在本项目中的主要流程已走完。

## 产品设计
 
```
 //目录结构
  src:
     |_controller //对接路由，媒体分发逻辑
     |_helpers    //存放项目中使用到的基本库
     |_model      //项目基础模型
     |_router     //路由
     |_main       //项目启动目录
 ```

## 细节规划
    待完善