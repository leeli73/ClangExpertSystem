/*
code50001@is not assigned a value
code50002@ninitialized variable
code50003@ossible null pointer dereference
code50004@which is out of bounds
code50005@is assigned a value that is never used
code50006@Memory leak
code50007@Resource leak
*/
error(code50001,"格式","变量为分配合法的值","为变量赋正确的值").
error(code50002,"错误","未初始化变量","正确初始化变量").
error(code50003,"警告","可能存在空指针","正确初始化指针").
error(code50004,"错误","存在数组越界","做好索引限制或扩容数组").
error(code50005,"格式","变量声明了却未使用","删除无用的变量").
error(code50006,"错误","内存泄露","及时释放无用的内存").
error(code50007,"错误","资源泄露","及时释放无用的读入资源").