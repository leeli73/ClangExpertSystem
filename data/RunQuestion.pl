/*
code10001 if语句的分支无法进入
code10002 while无法跳出
code10003 for无法跳出
code10004 switch多个分支的都被执行
code10005 获取的输入值不正确
code10006 程序运行时爆栈
code10007 变量总是不等
code10008 变量经计算,值却未变
code20001 使用了if
code20002 使用了switch&nbspcase
code20003 使用了判等
code20004 使用了scanf
code20005 使用了递归
code20006 使用了类似int&nbsparr[10],i;的语句
code20007 变量等于函数返回值
code20008 变量为引用或指针
code30001 使用的=而非==|$|a==1而非a=1
code30002 case中未使用break|$|case:code...;berak;
code30003 scanf未使用地址符|$|scanf(,&a);
code30004 递归深度太高|$|尝试换用迭代或其他思路
code30005 循环跳出条件出错|$|请检查循环跳出条件
code30006 函数返回值错误|$|请检查函数返回值的及相关运算
code30007 指针或引用使用错误|$|请检查指针的用法
code30008 栈溢出|$|请检查循环结束条件或拆分定义变量
*/
error(code10001,code20003,[code30001]).
error(code10002,code20003,[code30001]).
error(code10003,code20003,[code30001]).
error(code10004,code20002,[code30002]).
error(code10005,code20004,[code30003]).
error(code10006,code20005,[code30004]).
error(code10007,code20004,[code30003]).
error(code10008,code20007,[code30006]).
error(code10008,code20008,[code30007]).
error(code10002,code20006,[code30008]).
error(code10003,code20006,[code30008]).
