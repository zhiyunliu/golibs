# 高性能Go URL路径匹配器

一个支持参数化路径、通配符和多段匹配的高性能Go语言路径匹配器。

## 特性

- ✅ 参数化路径匹配：`/api/{servername}/userprofile/get`
- ✅ 单段通配符：`/api/user-asset/*/get`
- ✅ 多段通配符：`/api/user-asset/**`
- ✅ 高性能树形结构匹配
- ✅ 灵活的Option模式配置
- ✅ 路径分组支持
- ✅ 完整的验证和错误处理