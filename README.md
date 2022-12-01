# search
开源倒排搜索服务


# 目录结构
- core 倒排，索引存储，搜索
- api 提供接口，封装的实现
- http_service 提供http接口
- http_manage 提供http 管理界面接口
- tcp_service 提供tcp接口
- grpc_service 提供grpc接口
- utils 公共方法


# 概念

term : 索引词
lexicon :  索引词组，是一个树，前最树
inverted file : 存储倒排列表的数据文件
postingList : 倒排列表

# 排序

排序规则 ： 时序，排序值，词频
综合排序 ： 时序 > 词频
自定义排序: 排序值 > 时序


