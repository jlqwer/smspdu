# smspdu

### 获取发送指令

```
func StringToPDU(smsText string, phoneNumber string, smscNumber string, size int, mclass int, valid int, receipt bool) string  {}
```

以EC20为例：

```
package main

import (
   "fmt"
   "github.com/jlqwer/smspdu"
)

func main()  {
   cmd := smspdu.StringToPDU("这是短信内容", "8613800138000", "", 16, 2, 0, false)
   fmt.Println(cmd)
}
```

返回
```
AT+CMGS=26
0001000D91683108108300F0001A0C8FD9662F77ED4FE151855BB9
```