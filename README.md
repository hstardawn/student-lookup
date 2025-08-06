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

**接口地址：** `POST /api/v1/search`

**请求体：**
```json
{
    "student_id": "202301001",
    "name": "张三"
}
```

**响应示例：**
```json
{
    "success": true,
    "message": "查询成功",
    "data": {
        "year": "2023",
        "college": "计算机学院",
        "class": "计算机1班",
        "student_id": "202301001",
        "name": "张三",
        "major": "计算机科学与技术"
    }
}
```

### 2. 查询学生信息（GET）

**接口地址：** `GET /api/v1/search?student_id=202301001&name=张三`

**参数：**
- `student_id`: 学号（必需）
- `name`: 姓名（必需）

### 3. 重新加载Excel文件

**接口地址：** `POST /api/v1/reload`

当添加或修改Excel文件后，可调用此接口重新加载数据。

### 4. 获取统计信息

**接口地址：** `GET /api/v1/stats`

**响应示例：**
```json
{
    "success": true,
    "data": {
        "total_students": 1500,
        "colleges": ["计算机学院", "电子学院", "机械学院"],
        "years": ["2021", "2022", "2023"]
    }
}
```

### 5. 健康检查

**接口地址：** `GET /health`

## 使用示例

### 使用curl查询

```bash
# POST方式查询
curl -X POST http://localhost:8080/api/v1/search \
  -H "Content-Type: application/json" \
  -d '{"student_id":"202301001","name":"张三"}'

# GET方式查询
curl "http://localhost:8080/api/v1/search?student_id=202301001&name=张三"

# 获取统计信息
curl http://localhost:8080/api/v1/stats

# 重新加载数据
curl -X POST http://localhost:8080/api/v1/reload
```

### 使用JavaScript查询

```javascript
// 查询学生信息
async function searchStudent(studentId, name) {
    const response = await fetch('http://localhost:8080/api/v1/search', {
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
searchStudent('202301001', '张三').then(result => {
    if (result.success) {
        console.log('找到学生:', result.data);
    } else {
        console.log('查询失败:', result.message);
    }
});
```

## 错误处理

API会返回相应的HTTP状态码和错误信息：

- `200 OK`: 查询成功
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 未找到匹配的学生信息
- `500 Internal Server Error`: 服务器内部错误

## 注意事项

1. **Excel文件格式**：确保Excel文件包含正确的6列数据，第一行为标题行将被跳过
2. **数据完整性**：学号和姓名字段不能为空，否则该条记录会被跳过
3. **文件编码**：支持标准的xlsx格式文件
4. **查询匹配**：查询时学号和姓名必须完全匹配（区分大小写）
5. **性能考虑**：所有数据加载到内存中，适合中小型数据集

## 依赖包

- `github.com/gin-gonic/gin`: Web框架
- `github.com/xuri/excelize/v2`: Excel文件处理

## License

MIT License
