### 编译

```
git clone ...
cd sample-code-golang
export GOPATH=$PWD
make
```

### 用法

各测试后加-h或--help可显示出其用法与各命令参数。
如

```
./wstts -h
Usage of ./wstts:
  -auth
    	Need auth?
  -codec string
    	Codec (default "opu2")
  -count int
    	Test count (default 1)
  -file string
    	Out file
  -host string
    	Server address (default "wss://apigwws.open.rokid.com/api")
  -lang string
    	Language (default "zh")
  -text string
    	Tts Text (default "今天天气怎么样?")
```
