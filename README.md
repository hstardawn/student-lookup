# 学生信息查询API

这是一个基于Go语言开发的HTTP API服务，用于根据学号和姓名查找Excel文件中的学生信息。

## 功能特性

- 🔍 根据学号和姓名精确查找学生信息
- 📊 支持多个Excel文件同时加载
- 🔄 支持热重载Excel文件
- 📈 提供数据统计信息
- 🌐 RESTful API设计
- ✅ 健康检查接口

## Excel文件格式

Excel文件应包含以下6列（按顺序）：
1. **年份** - 入学年份
2. **学院名称** - 所属学院
3. **班级** - 班级信息
4. **学号** - 学生学号（用于查询）
5. **姓名** - 学生姓名（用于查询）
6. **录取专业名称** - 专业信息

## 项目结构

```
student-lookup/
├── main.go           # 主程序文件
├── go.mod           # Go模块文件
├── data/            # Excel文件存放目录
│   └── *.xlsx       # 学生信息Excel文件
└── README.md        # 项目说明文档
```

## 安装和运行

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 准备数据

将包含学生信息的Excel文件（.xlsx格式）放入 `data/` 目录下。

### 3. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API接口

### 1. 查询学生信息（POST）

**接口地址：** `POST /api/search`

**请求体：**
```json
{
    "student_id": "302024670001",
    "name": "张伟"
}
```

**响应示例：**

成功找到学生：
```json
{
    "code": 200,
    "msg": "查询成功",
    "data": "计算机科学与技术"
}
```

未找到学生：
```json
{
    "code": 200001,
    "msg": "非新生"
}
```

系统错误：
```json
{
    "code": 200500,
    "msg": "学号格式错误，应为12位数字"
}
```

### 2. 查询学生信息（GET）

**接口地址：** `GET /api/search?student_id=302024670001&name=张伟`

**参数：**
- `student_id`: 学号（必需，12位数字）
- `name`: 姓名（必需）

**响应格式同POST接口**

### 3. 健康检查

**接口地址：** `GET /health`

**响应示例：**
```json
{
    "status": "ok",
    "total_records": 38
}
```

## 响应状态码说明

| Code | 含义 | 说明 |
|------|------|------|
| 200 | 成功 | 查询成功，找到学生信息 |
| 200001 | 非新生 | 查询成功，但未找到匹配的学生信息 |
| 200500 | 系统错误 | 参数错误、学号格式错误等系统异常 |

注意：所有接口都返回HTTP 200状态码，具体的业务状态通过response中的code字段区分。

## 使用示例

### 使用curl查询

```bash
# POST方式查询
curl -X POST http://localhost:8080/api/search \
  -H "Content-Type: application/json" \
  -d '{"student_id":"302024670001","name":"张伟"}'

# GET方式查询
curl "http://localhost:8080/api/search?student_id=302024670001&name=张伟"

# 健康检查
curl http://localhost:8080/health
```

### 使用JavaScript查询

```javascript
// 查询学生信息
async function searchStudent(studentId, name) {
    const response = await fetch('http://localhost:8080/api/search', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            student_id: studentId,
            name: name
        })
    });
    
    const result = await response.json();
    return result;
}

// 使用示例
searchStudent('302024670001', '张伟').then(result => {
    if (result.code === 200) {
        console.log('找到学生，专业:', result.data);
    } else if (result.code === 200001) {
        console.log('非新生');
    } else {
        console.log('查询错误:', result.msg);
    }
});
```

## 注意事项

1. **Excel文件格式**：确保Excel文件包含正确的6列数据，第一行为标题行将被跳过
2. **数据完整性**：学号和姓名字段不能为空，否则该条记录会被跳过
3. **学号格式**：学号必须是12位数字，如：302024670001
4. **文件编码**：支持标准的xlsx格式文件
5. **查询匹配**：查询时学号和姓名必须完全匹配（区分大小写）
6. **响应格式**：所有接口统一返回HTTP 200，通过code字段区分业务状态
7. **性能考虑**：所有数据加载到内存中，适合中小型数据集

## 依赖包

- `github.com/gin-gonic/gin`: Web框架
- `github.com/xuri/excelize/v2`: Excel文件处理

## License

MIT License
