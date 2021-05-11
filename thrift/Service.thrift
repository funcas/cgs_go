// thrift -r --gen go:package_prefix=github.com/funcas/cgs/gen-go/ thrift/Service.thrift

namespace go process

union ObjectArg {
  1: i32 int_arg;
  2: i64 long_arg;
  3: string string_arg;
  4: bool bool_arg;
  5: binary binary_arg;
  6: double double_arg;
}

typedef map<string, ObjectArg> data
typedef list<data> cudData

struct Resp {
    1: string transCode;
    2: data data;
    3: cudData cudData;
}


//定义服务
service EntryService {
    Resp Execute(
        1: string transCode
    )

    Resp ExecuteWithParams(
        1: string transCode
        2: data params
    )
}