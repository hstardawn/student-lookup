package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// Student 学生信息结构体
type Student struct {
	Year      string `json:"year"`       // 年份
	College   string `json:"college"`    // 学院名称
	Class     string `json:"class"`      // 班级
	StudentID string `json:"student_id"` // 学号
	Name      string `json:"name"`       // 姓名
	Major     string `json:"major"`      // 录取专业名称
}

// SearchRequest 查询请求结构体
type SearchRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

// SearchResponse 查询响应结构体
type SearchResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data,omitempty"`
}

var students []Student

// loadExcelFiles 加载所有Excel文件
func loadExcelFiles() error {
	// 数据目录路径
	dataDir := "./data"

	// 获取所有xlsx文件
	files, err := filepath.Glob(filepath.Join(dataDir, "*.xlsx"))
	if err != nil {
		return fmt.Errorf("获取Excel文件失败: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("在 %s 目录下未找到xlsx文件", dataDir)
	}

	students = []Student{} // 清空现有数据

	for _, file := range files {
		if err := loadExcelFile(file); err != nil {
			log.Printf("加载文件 %s 失败: %v", file, err)
			continue
		}
		log.Printf("成功加载文件: %s", file)
	}

	log.Printf("总共加载了 %d 条学生记录", len(students))
	return nil
}

// loadExcelFile 加载单个Excel文件
func loadExcelFile(filename string) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取工作表，优先选择"学生信息"工作表
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("Excel文件中没有工作表")
	}

	var sheetName string
	// 优先选择"学生信息"工作表
	for _, sheet := range sheets {
		if sheet == "学生信息" {
			sheetName = sheet
			break
		}
	}
	// 如果没有找到"学生信息"工作表，使用第一个工作表
	if sheetName == "" {
		sheetName = sheets[0]
	}

	log.Printf("文件 %s 选择工作表: %s", filename, sheetName)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("读取工作表数据失败: %v", err)
	}

	// 跳过标题行，从第二行开始读取
	for i, row := range rows {
		if i == 0 {
			continue // 跳过标题行
		}

		// 确保行有足够的列
		if len(row) < 6 {
			log.Printf("文件 %s 第 %d 行数据不完整，跳过", filename, i+1)
			continue
		}

		student := Student{
			Year:      strings.TrimSpace(row[0]),
			College:   strings.TrimSpace(row[1]),
			Class:     strings.TrimSpace(row[2]),
			StudentID: strings.TrimSpace(row[3]),
			Name:      strings.TrimSpace(row[4]),
			Major:     strings.TrimSpace(row[5]),
		}

		// 验证必要字段不为空
		if student.StudentID != "" && student.Name != "" {
			students = append(students, student)
		}
	}

	return nil
}

// isValidStudentID 验证学号格式（12位数字）
func isValidStudentID(studentID string) bool {
	// 检查长度是否为12位
	if len(studentID) != 12 {
		return false
	}

	// 检查是否都是数字
	for _, char := range studentID {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// searchStudent 查找学生信息
func searchStudent(studentID, name string) *Student {
	for _, student := range students {
		if student.StudentID == studentID && student.Name == name {
			return &student
		}
	}
	return nil
}

// handleSearch 处理查询请求
func handleSearch(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, SearchResponse{
			Code: 200500,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证学号格式
	if !isValidStudentID(req.StudentID) {
		c.JSON(http.StatusOK, SearchResponse{
			Code: 200500,
			Msg:  "学号格式错误，应为12位数字",
		})
		return
	}

	// 查找学生
	student := searchStudent(req.StudentID, req.Name)
	if student == nil {
		c.JSON(http.StatusOK, SearchResponse{
			Code: 200001,
			Msg:  "非新生",
		})
		return
	}

	c.JSON(http.StatusOK, SearchResponse{
		Code: 200,
		Msg:  "查询成功",
		Data: student.Major,
	})
}

// handleSearchByParams 通过URL参数查询
func handleSearchByParams(c *gin.Context) {
	studentID := c.Query("student_id")
	name := c.Query("name")

	if studentID == "" || name == "" {
		c.JSON(http.StatusOK, SearchResponse{
			Code: 200500,
			Msg:  "缺少必要参数: student_id 和 name",
		})
		return
	}

	// 验证学号格式
	if !isValidStudentID(studentID) {
		c.JSON(http.StatusOK, SearchResponse{
			Code: 200500,
			Msg:  "学号格式错误，应为12位数字",
		})
		return
	}

	// 查找学生
	student := searchStudent(studentID, name)
	if student == nil {
		c.JSON(http.StatusOK, SearchResponse{
			Code: 200001,
			Msg:  "非新生",
		})
		return
	}

	c.JSON(http.StatusOK, SearchResponse{
		Code: 200,
		Msg:  "查询成功",
		Data: student.Major,
	})
}

func main() {
	// 加载Excel文件
	if err := loadExcelFiles(); err != nil {
		log.Fatalf("加载Excel文件失败: %v", err)
	}

	// 创建Gin路由
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// API路由
	api := r.Group("/api")
	{
		api.POST("/search", handleSearch)        // POST方式查询
		api.GET("/search", handleSearchByParams) // GET方式查询
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":        "ok",
			"total_records": len(students),
		})
	})

	// 根路径
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "学生信息查询API",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"POST /api/search": "查询学生信息 (JSON)",
				"GET /api/search":  "查询学生信息 (URL参数)",
				"GET /health":      "健康检查",
			},
		})
	})

	port := ":8080"
	log.Printf("服务器启动在端口 %s", port)
	log.Printf("访问 http://localhost%s 查看API信息", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
