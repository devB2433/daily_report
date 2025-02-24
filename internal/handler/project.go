package handler

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"daily-report/internal/database"
	"daily-report/internal/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GetProjects 获取项目列表
func GetProjects(c *gin.Context) {
	db := database.GetDB()
	var projects []model.Project

	// 查询所有项目，按状态和名称排序
	if err := db.Order("FIELD(status, 'active', 'suspended', 'completed'), name ASC").Find(&projects).Error; err != nil {
		// log.Printf("Failed to get projects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取项目列表失败",
		})
		return
	}

	// 如果没有找到任何项目
	if len(projects) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []model.Project{},
			"message": "暂无项目",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projects,
	})
}

// CreateProject 创建新项目
func CreateProject(c *gin.Context) {
	// log.Printf("开始处理创建项目请求")

	var project model.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		// log.Printf("解析项目数据失败: %v, 请求体: %+v", err, c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	// log.Printf("接收到的项目数据: %+v", project)

	// 验证必填字段
	if project.Name == "" || project.Code == "" {
		// log.Printf("必填字段验证失败 - name: %s, code: %s", project.Name, project.Code)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目名称和代号不能为空",
		})
		return
	}

	// 验证项目状态
	validStatus := map[string]bool{
		"active":    true,
		"completed": true,
		"suspended": true,
	}
	if !validStatus[project.Status] {
		// log.Printf("项目状态无效: %s", project.Status)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目状态无效",
		})
		return
	}

	// log.Printf("开始数据库操作")
	db := database.GetDB()

	// 检查项目名称是否已存在
	var existingProject model.Project
	if err := db.Where("name = ?", project.Name).First(&existingProject).Error; err == nil {
		// log.Printf("项目名称已存在: %s", project.Name)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目名称已存在",
		})
		return
	}

	// 检查项目代号是否已存在
	if err := db.Where("code = ?", project.Code).First(&existingProject).Error; err == nil {
		// log.Printf("项目代号已存在: %s", project.Code)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "项目代号已存在",
		})
		return
	}

	// 创建项目
	if err := db.Create(&project).Error; err != nil {
		// log.Printf("创建项目失败: %+v, 错误: %v", project, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建项目失败",
		})
		return
	}

	// log.Printf("项目创建成功: ID=%d, Name=%s", project.ID, project.Name)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建项目成功",
		"data":    project,
	})
}

// UpdateProject 更新项目信息
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	var project model.Project

	db := database.GetDB()
	if err := db.First(&project, id).Error; err != nil {
		log.Printf("Failed to find project: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "项目不存在",
		})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据无效",
		})
		return
	}

	if err := db.Save(&project).Error; err != nil {
		log.Printf("Failed to update project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新项目失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新项目成功",
		"data":    project,
	})
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	// 检查是否存在与该项目相关联的日报
	var count int64
	db.Model(&model.Task{}).Where("project_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无法删除项目，因为存在与该项目相关联的日报",
		})
		return
	}

	if err := db.Delete(&model.Project{}, id).Error; err != nil {
		log.Printf("Failed to delete project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除项目失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除项目成功",
	})
}

// GetProject 获取项目详情
func GetProject(c *gin.Context) {
	id := c.Param("id")
	log.Printf("获取项目详情请求，项目ID: %s", id)

	db := database.GetDB()

	var project model.Project
	if err := db.First(&project, id).Error; err != nil {
		log.Printf("查找项目失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("项目不存在: %v", err),
		})
		return
	}

	log.Printf("成功获取项目详情: ID=%d, Name=%s", project.ID, project.Name)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    project,
	})
}

// ExportProjects 导出项目数据为CSV格式
func ExportProjects(c *gin.Context) {
	db := database.GetDB()
	var projects []model.Project

	if err := db.Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取项目数据失败",
		})
		return
	}

	// 设置响应头
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=projects.csv")

	// 写入UTF-8 BOM
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// 创建CSV写入器
	writer := csv.NewWriter(c.Writer)

	// 写入表头
	writer.Write([]string{"code", "name", "status", "manager", "description"})

	// 写入项目数据
	for _, project := range projects {
		writer.Write([]string{
			project.Code,
			project.Name,
			project.Status,
			project.Manager,
			project.Description,
		})
	}

	writer.Flush()
}

// ImportProjects 从CSV文件导入项目数据
func ImportProjects(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("文件上传错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请选择要导入的文件",
		})
		return
	}

	log.Printf("接收到文件: %s, 大小: %d bytes", file.Filename, file.Size)

	// 检查文件类型
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请上传CSV文件",
		})
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		log.Printf("打开文件错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "无法读取文件",
		})
		return
	}
	defer src.Close()

	// 读取前三个字节以检查 BOM
	bom := make([]byte, 3)
	_, err = src.Read(bom)
	if err != nil {
		log.Printf("读取 BOM 错误: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "读取文件失败",
		})
		return
	}

	// 如果文件以 UTF-8 BOM 开头，跳过这些字节
	if bom[0] == 0xEF && bom[1] == 0xBB && bom[2] == 0xBF {
		log.Printf("检测到 UTF-8 BOM")
	} else {
		// 如果不是 BOM，将文件指针重置到开始位置
		_, err = src.Seek(0, 0)
		if err != nil {
			log.Printf("重置文件指针错误: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "读取文件失败",
			})
			return
		}
	}

	// 创建 CSV 读取器
	reader := csv.NewReader(src)
	reader.LazyQuotes = true       // 允许字段中的引号
	reader.FieldsPerRecord = -1    // 允许变长记录
	reader.TrimLeadingSpace = true // 去除前导空格

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("CSV读取错误: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("CSV文件格式无效: %v", err),
		})
		return
	}

	// 检查并修正编码问题
	for i := range records {
		for j := range records[i] {
			// 检查是否包含无效的UTF-8字符
			if !utf8.ValidString(records[i][j]) {
				// 尝试将可能的GBK编码转换为UTF-8
				data, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(records[i][j]), simplifiedchinese.GBK.NewDecoder()))
				if err == nil {
					records[i][j] = string(data)
				}
			}
			// 清理不可见字符
			records[i][j] = strings.Map(func(r rune) rune {
				if unicode.IsPrint(r) || unicode.IsSpace(r) {
					return r
				}
				return -1
			}, records[i][j])
		}
	}

	if len(records) < 2 {
		log.Printf("CSV文件为空或只有表头: 行数=%d", len(records))
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "CSV文件不能为空",
		})
		return
	}

	// 验证表头
	header := records[0]
	expectedHeader := []string{"code", "name", "status", "manager", "description"}
	log.Printf("实际表头: %v", header)
	log.Printf("期望表头: %v", expectedHeader)

	if len(header) != len(expectedHeader) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("表头列数不正确，期望 %d 列，实际 %d 列", len(expectedHeader), len(header)),
		})
		return
	}

	if !reflect.DeepEqual(header, expectedHeader) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("CSV文件格式不正确，请使用正确的模板。\n期望的表头: %v\n实际的表头: %v", expectedHeader, header),
		})
		return
	}

	db := database.GetDB()
	var successCount, errorCount int
	var errorMessages []string

	// 开始事务
	tx := db.Begin()

	// 处理每一行数据
	for i, record := range records[1:] {
		lineNum := i + 2 // 实际行号（跳过表头，从第2行开始）
		if len(record) != 5 {
			errMsg := fmt.Sprintf("第 %d 行: 列数不正确，期望 5 列，实际 %d 列", lineNum, len(record))
			errorMessages = append(errorMessages, errMsg)
			log.Printf(errMsg)
			errorCount++
			continue
		}

		code := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		status := strings.TrimSpace(record[2])
		manager := strings.TrimSpace(record[3])
		description := strings.TrimSpace(record[4])

		// 验证必填字段
		if code == "" || name == "" {
			errMsg := fmt.Sprintf("第 %d 行: 项目编号和名称不能为空 (code=%s, name=%s)", lineNum, code, name)
			errorMessages = append(errorMessages, errMsg)
			log.Printf(errMsg)
			errorCount++
			continue
		}

		// 验证状态值
		if status != "active" && status != "completed" && status != "suspended" {
			errMsg := fmt.Sprintf("第 %d 行: 状态值无效 '%s'，必须是 active、completed 或 suspended", lineNum, status)
			errorMessages = append(errorMessages, errMsg)
			log.Printf(errMsg)
			errorCount++
			continue
		}

		// 查找是否存在相同编号的项目
		var existingProject model.Project
		if err := tx.Where("code = ?", code).First(&existingProject).Error; err == nil {
			// 更新现有项目，只更新CSV中包含的字段
			updates := map[string]interface{}{
				"name":        name,
				"status":      status,
				"manager":     manager,
				"description": description,
			}
			if err := tx.Model(&existingProject).Updates(updates).Error; err != nil {
				errMsg := fmt.Sprintf("第 %d 行: 更新项目失败 '%s': %v", lineNum, code, err)
				errorMessages = append(errorMessages, errMsg)
				log.Printf(errMsg)
				errorCount++
				continue
			}
			log.Printf("第 %d 行: 更新项目成功: %s", lineNum, code)
		} else {
			// 创建新项目
			project := model.Project{
				Code:        code,
				Name:        name,
				Status:      status,
				Manager:     manager,
				Description: description,
			}
			if err := tx.Create(&project).Error; err != nil {
				errMsg := fmt.Sprintf("第 %d 行: 创建项目失败 '%s': %v", lineNum, code, err)
				errorMessages = append(errorMessages, errMsg)
				log.Printf(errMsg)
				errorCount++
				continue
			}
			log.Printf("第 %d 行: 创建项目成功: %s", lineNum, code)
		}
		successCount++
	}

	// 提交或回滚事务
	if errorCount > 0 {
		tx.Rollback()
		errorSummary := strings.Join(errorMessages, "\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("导入失败：发现 %d 个错误:\n%s", errorCount, errorSummary),
		})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("成功导入 %d 条记录", successCount),
	})
}
