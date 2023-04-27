# license_check 使用说明

## 1.主要功能

#### license_check包主要用于验证用户的license的有效性，是否过期等

## 2.使用示例

#### 如demo所示，用户需要提供一个字符串类型的公钥，并在license_check包下加入license文件。然后调用CheckLicense接口，便会启动服务。CheckLicense每隔大约30秒便会对license进行一次检验，检验其是否过期。控制台会显示输出一条对应的日志信息，显示license有效性检验结果。

## 3.输入

#### 参数1：公钥 publickey 类型：string

#### 参数2：许可 license 类型：文本文件

## 4.输出

#### 1) 日志配置信息：

![log](./assets/log.png)

#### 2) 如果license验证通过，将会看到如下信息：

![log](./assets/pass.png)

#### 3) 当license验证不通过时，控制台会显示验证失败信息，并中止程序

![log](./assets/fail.png)

## 5.注意事项

#### 1）可以通过修改日志配置文件log.json改变日志输出等级，输出格式等输出信息

#### 2）当启动license_check服务时，会在当前目录下自动生成一个app.log文件保存日志信息