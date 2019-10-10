# XSQL3

数据库深度优化，颠覆性操作习惯，数据表自动生成，操作完全基于数据对象，减少因为拼写错误导致的错误，减少对sql语句的记忆。***献给我心爱姑娘的礼物***


第一步：在config.xml里面配置数据库属性

```xml
<?xml version ="1.0" encoding="utf-8"?>
<config ver="1.0.0.0">
  <db_info echo="mysql数据库配置">
    <db_server dynamic="no" echo="数据库地址">192.168.1.145</db_server>
    <db_port dynamic="no" echo="数据库端口">3306</db_port>
    <db_name dynamic="no" echo="数据库名称">mbo_server_db</db_name>
    <db_user dynamic="no" echo="数据库用户名">root</db_user>
    <db_password dynamic="no" echo="数据库密码">123456</db_password>
    <db_charset dynamic="no" echo="数据库客户端字符集">utf8mb4</db_charset>
  </db_info>
</config>
```

第二步：定义数据表

```go
//基本信息;
type BridgeBaseBasic struct {
  ID        Int    `json:"id" type:"int(11)" comment:"桥梁信息id" required:"yes" mark:"NOT NULL PRIMARY KEY AUTO_INCREMENT"`
  BridgeId  Int    `json:"bridgeId" type:"int(11)" default:"0" comment:"桥梁信息id" required:"no"`
  Len       Float  `json:"len" type:"float(11)" default:"0" comment:"桥梁信息id" required:"no"`
  BseType   Int    `json:"bseType" type:"tinyint(3)" Q:"" default:"0" comment:"类型，0 QM/桥面；1 ZJ/桩基；2 CT/承台；3 XL/系梁；4 GL/盖梁；5 DZ/墩柱；6 TS/台身；7 XJL/现浇段；8 YZL/预制梁；9 GXL/钢箱梁"`
  BseUnitId String `json:"bseUnitId" type:"varchar(60)" Q:"" default:"" comment:"1编号" required:"yes"`
  Note      String `json:"note" type:"varchar(255)" Q:"" default:"" comment:"备注" required:"yes"`
}

func (t *BridgeBaseBasic) TableName() string {
   return TestSql
}

func (t *BridgeBaseBasic) NewInterface() DBOperation {
   return xsql3.ScanStructInterface(t).(*BridgeBaseBasic)
}

func NewInterface() *BridgeBaseBasic {
   return (&BridgeBaseBasic{}).NewInterface().(*BridgeBaseBasic)
}

func init() {
  f := xsql3.Register("test_sql",NewInterface(),"测试数据表")
  f.AddUniqueKey("mbo_sql_Unique",[]string{"id"})
  f.AddIndexKey("mbo_sql_index",[]string{"BridgeId"})
  f.CopyToMySQL()
}
```

1、该处使用的Int、String、Float是github.com/NiuStar/xsql3/Type中定义的数据类型，必须使用这种数据类型

2、ScanStructInterface为github.com/NiuStar/xsql3中的一个方法，该方法为对结构体对象初始化作用

3、init在外部引用时会被调用，开始注册这个表，会在数据库中生成该表，并针对结构的修改，都会实时更新在数据表中

第三步开始使用：

所有的数据库相关操作都是通过DBFactory.NewDBFactory()来打交道。

1、数据插入：

方法名：InsertDB

参数：表对象，通过第二步定义生成的对象

返回值：本次插入的id

```go
  bridge := DB.NewInterface()
  bridge.BseUnitId.SetValue("BseUnitId——哈哈")
  bridge.Node.SetValue("node20000-f返回")
  fmt.Println("插入的id ： ",DBFactory.NewDBFactory().InsertDB(bridge))
```

2、数据更新

方法名：UpdateDB

参数：表对象，如需更新的字段，不要为空，不需要更新的字段为空即可

返回值：本次更新影响的条数

```go
bridge := DB.NewInterface()
bridge.BseUnitId.SetValue("BseUnitId")
fmt.Println("影响条数 ： ",DBFactory.NewDBFactory().SetConditions(bridge.ID.Equal(1)).UpdateDB(bridge))
```

3、删除内容

方法名：UpdateDB

返回值：本次删除影响的条数

```go
test := DB.NewTestInterface()
fmt.Println("影响条数 ：",DBFactory.NewDBFactory().SetConditions(test.ID.Equal(2)).DeleteDB(test))
```

4、一个最简单的查询，如查询一个id为1的bridge这个对象的BseUnitId成员变量的内容。

获取查询结构有两种方式：
GetResultsOperation()与GetResults()

其中GetResultsOperation返回值会转化为对应的对象，而GetResults没有经过对象转换，就是数据库查出来的内容

```go
bridge := DB.NewInterface()
list := DBFactory.NewDBFactory().SetTable(bridge).
    SetFields(&bridge.BseUnitId).
      SetConditions(bridge.ID.Equal(1)).
    GetResultsOperation()
```

5、复杂查询

```go
bridge := DB.NewInterface()
test := DB.NewTestInterface()
list := f.SetTable(bridge).
   SetFields(&bridge.BseUnitId).
      SetJoins(j).
      UNIONDISTINCTDB(f1).
      Count().
      SetFields(&bridge.ID,&bridge.BseUnitId,&bridge.Len,&bridge.ID,&test.ID,&test.Name).
      SetConditions(test.ID.EqualFold(&bridge.ID).AddCondition(test.ID.GreaterEqual(1))).
      SetORConditions(test.ID.EqualFold(&bridge.ID).AddCondition(test.ID.GreaterEqual(1))).
      AddORCondition(test.ID.EqualFold(&bridge.ID).AddCondition(test.ID.GreaterEqual(2))).
      SetGroupBy(&bridge.ID).
      SetOrderBy(&bridge.ID).
      DESC().
      Limit(0,10).
   GetResultsOperation()
```



[更多函数](https://github.com/NiuStar/xsql3/blob/master/DBFactory.md)