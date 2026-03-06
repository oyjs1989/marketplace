# 基于项目代码提炼的编码规范规则

## 基于项目代码分析提炼的新规则

1. **结构体命名采用驼峰命名法，首字母大写**
   - 所有结构体名称都使用驼峰命名法，首字母大写，如`ProviderManager`、`CloudMgrServiceImpl`等

2. **结构体字段命名采用驼峰命名法，首字母大写**
   - 结构体字段名使用驼峰命名法，首字母大写，如`PolicyID`、`ProviderCode`等

3. **函数和方法命名采用驼峰命名法，首字母小写（私有）或大写（公共）**
   - 私有函数/方法：首字母小写，如`getActions`、`createRole`
   - 公共函数/方法：首字母大写，如`GetActions`、`CreateRole`

4. **函数返回值中error类型永远放在最后**
   - 所有函数签名中，error类型作为最后一个返回值，如`func (m *ProviderManager) GetActions(ctx context.Context, product string) ([]*cloud_api.Action, error)`

5. **错误处理必须使用errors包而不是fmt.Errorf**
   - 项目中统一使用`errors.New()`、`errors.Wrap()`、`errors.Wrapf()`等函数处理错误
   - 避免使用`fmt.Errorf()`构造错误信息

6. **错误信息应使用英文描述**
   - 错误信息统一使用英文描述，如`return errors.New("role alias cannot be empty")`

7. **日志记录使用结构化日志字段**
   - 日志记录时使用结构化字段，如`log.String("alias", request.GetAlias())`、`log.Int("action_count", len(addActionList))`
   - 避免使用字符串拼接方式记录日志

8. **日志级别合理使用**
   - 错误信息使用`log.Error()`
   - 调试信息使用`log.Debug()`
   - 普通信息使用`log.Info()`

9. **常量统一定义在constants包中**
   - 所有常量都定义在`internal/app/constants/`目录下对应的文件中
   - 常量命名使用大写字母和下划线，如`RoleDescription`、`RoleSessionDuration`

10. **业务对象（BO）与数据模型（Model）分离**
    - 业务对象定义在`internal/app/biz/bo/`目录下
    - 数据库模型定义在`internal/app/model/db/`目录下

11. **接口定义与实现分离**
    - 接口定义通常在单独的文件中，如`internal/app/biz/repo/`目录下的接口定义
    - 接口实现在对应的impl文件中，如`internal/app/data/db/`目录下的实现

12. **函数参数较多时使用结构体封装**
    - 当函数参数较多时，通过结构体封装参数传递，避免长参数列表

13. **上下文（context）作为函数的第一个参数**
    - 所有需要上下文的函数，ctx作为第一个参数传入

14. **返回错误时立即处理，避免延迟处理**
    - 错误检查紧跟在可能产生错误的函数调用之后，如`if err != nil { return err }`

15. **注释格式统一，公共函数必须有注释**
    - 公共函数/方法必须有注释说明功能、参数和返回值
    - 注释格式遵循统一标准

16. **包名与目录名保持一致**
    - Go包名与目录名保持一致，如`internal/app/biz/param_template.go`中包名为`biz`

17. **包导入按标准库、第三方库、内部库分组**
    - 包导入按标准库、第三方库、内部库顺序分组，每组之间空一行

18. **测试文件命名规范**
    - 测试文件以`_test.go`结尾，与被测试文件保持一致的命名

19. **测试函数命名规范**
    - 测试函数以`Test`开头，后跟被测试函数名

20. **业务逻辑与数据访问分离**
    - 业务逻辑在`internal/app/biz/`目录下实现
    - 数据访问在`internal/app/data/`目录下实现

21. **云厂商相关代码按厂商分目录组织**
    - 不同云厂商的实现代码按厂商分目录，如`internal/pkg/cloud_api/aliyun/`、`internal/pkg/cloud_api/aws/`等

22. **配置相关代码统一放在config目录**
    - 配置管理相关代码放在`internal/app/config/`目录下

23. **工具类代码放在pkg目录**
    - 公共工具类代码放在`internal/pkg/`目录下

24. **错误日志必须包含错误字段**
    - 记录错误日志时必须使用`log.ErrorField(err)`包含错误信息

25. **函数文档注释格式统一**
    - 公共函数使用统一的文档注释格式，说明函数功能、参数和返回值

26. **变量命名具有描述性**
    - 变量命名具有描述性，避免使用单字母或无意义的缩写

27. **避免魔法数字和字符串**
    - 避免在代码中直接使用魔法数字和字符串，使用常量或配置

28. **函数长度控制**
    - 函数长度适中，过长的函数应拆分为多个小函数

29. **依赖注入原则**
    - 通过构造函数注入依赖，避免在函数内部直接创建依赖对象

30. **接口设计遵循最小原则**
    - 接口设计遵循最小原则，只包含必要的方法

31. **错误处理统一包装**
    - 错误处理统一使用errors包进行包装，保留错误堆栈信息

32. **日志上下文传递**
    - 通过`log.With()`传递日志上下文信息，确保日志关联性

33. **结构体嵌入用于扩展**
    - 使用结构体嵌入实现代码复用和扩展

34. **初始化函数集中管理**
    - 初始化相关代码集中在`internal/app/server/server.go`中管理

35. **服务实现统一注册**
    - 服务实现在`internal/app/service/`目录下统一管理