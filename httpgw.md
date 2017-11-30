## 流程
```
client                   httpgw                   InternalService(grpc)
-------------------------------------------------------------------------

http(带Authorization)      -->   根据Authorization进行设备认证或用户认证
                        A  <--   认证失败 500 body为失败原因的字符串
                           或
                           -->   按url推断出InternalService的域名
                                 将body解析成req
                                 填充AccountId, DeviceTypeId, DeviceId(设备认证)或ClientId, UserId(用户认证)
                                 调用resp = grpc(req)
                        B  <--   调用失败 500 body为失败原因的字符串
                           或
                        C  <--   返回200 body为resp
```
client使用open/vN/packageName/serviceName.proto
InternalService使用inner/vN/packageName/serviceName.proto

## 请求url组成
* https://httpgw.open.rokid.com/vN/packageName/serviceName/methodName
* https://域名/版本/包名/服务名/方法名

## 设备认证
* Authorization: version={version};time={time};sign={sign};key={key};device_type_id={device_type_id};device_id={device_id};service={service}

sign的生成加密方式：

key={key}&device_type_id={device_type_id}&device_id={device_id}&service={service}&version={version}&time={time}&secret={secret}

的utf8字符串的md5值

其中{xxx}由xxx的值替代

key及secret由开发方通过管理平台获取，并保管。

## 用户认证
* Authorization: version={version};sign={sign};token={token};client_id={client_id};time={time}

sign的生成加密方式：

version={version}&token={token}&client_id={client_id}&time={time}&secret={secret}

的utf8字符串的md5值

其中{xxx}由xxx的值替代

client_id及secret由用户通过平台获取，并保管。
