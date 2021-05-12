// thrift -r --gen go:package_prefix=github.com/funcas/cgs/gen-go/ thrift/Service.thrift

namespace go process

typedef map<string, string> data
typedef list<data> cudData

struct Resp {
    1: string transCode;
    2: string data;   // 无解析，返回原始报文
    3: data uniData;  // 解析过的非数据集数据
    4: map<string,cudData> cudData; //解析出来的数据集数据, key为关联键
}


service EntryService {
    // 执行不带参数的交易
    Resp execute(1: string transCode)

    // 执行带参数的交易
    Resp executeWithParams(1: string transCode, 2: map<string,string> params)

    // 重载配置文件
    Resp reload()

}