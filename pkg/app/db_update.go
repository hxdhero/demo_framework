package app

import (
	"fmt"
	"lls_api/pkg/rerr"
	"lls_api/pkg/util/text"
	"reflect"
	"strings"
)

// Update 根据模型和字段更新
// user := model.BlueUserSaas{ID:1,Name:"zhangsan"}
// ctx.DB().Update(&user,"name")
// model 需要跟新的实例
// columns 需要更新的字段
func (d *Database) Update(model Moduler, columns ...string) Result {
	result := Result{}

	if model == nil {
		result.Error = rerr.New("model cannot be nil")
		return result
	}
	if model.GetID() <= 0 {
		result.Error = rerr.New("模型的id错误")
		return result
	}

	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr {
		modelValue = modelValue.Elem()
	}

	if !modelValue.IsValid() {
		result.Error = rerr.New("invalid model value")
		return result
	}

	modelType := modelValue.Type()

	// 创建 tag 到字段名的映射
	tagToField := make(map[string]reflect.StructField)
	fieldValueMap := make(map[string]reflect.Value) // 存储传入模型的字段值

	// 遍历所有字段，建立映射关系（包括嵌套结构体的字段）
	buildFieldMapping(modelValue, modelType, "", &tagToField, &fieldValueMap)

	var updateInfos []string
	// 根据传入的 column 名称查找对应字段
	for _, columnName := range columns {
		_, exists := tagToField[columnName]
		if !exists {
			result.Error = rerr.Errorf("字段未找到: %s", columnName)
			return result
		}

		fieldValue, valueExists := fieldValueMap[columnName]
		var updateValue any
		if valueExists && fieldValue.IsValid() {
			updateValue = fieldValue.Interface()
		}

		updateInfos = append(updateInfos, fmt.Sprintf("%s=%v", columnName, updateValue))
	}

	res := d.db().Select(columns).Updates(model)
	d.Log().InfofSkip(1, "更新表:%s, 数据:%s", model.TableName(), strings.Join(updateInfos, ","))
	result.Error = rerr.Wrap(res.Error)
	result.RowsAffected = res.RowsAffected
	return result
}

// 递归构建字段映射，处理嵌套结构体
func buildFieldMapping(modelValue reflect.Value, modelType reflect.Type, prefix string, tagToField *map[string]reflect.StructField, fieldValueMap *map[string]reflect.Value) {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		fieldValue := modelValue.Field(i)

		// 跳过不可导出的字段
		if !fieldValue.CanInterface() {
			continue
		}

		// 获取 gorm tag 中的 column 名称
		gormTag := field.Tag.Get("gorm")
		columnName := getColumnFromGormTag(gormTag)

		// 如果没有指定 column，则使用字段名的 snake_case
		if columnName == "" {
			if prefix != "" {
				columnName = prefix + "_" + text.ToSnakeCase(field.Name)
			} else {
				columnName = text.ToSnakeCase(field.Name)
			}
		} else {
			// 如果有前缀，需要处理 column name
			if prefix != "" {
				// 如果 gorm tag 中指定了 column，保持原样；否则添加前缀
				if !strings.Contains(gormTag, "column:") {
					columnName = prefix + "_" + columnName
				}
			}
		}

		// 检查字段类型是否为结构体或指针指向结构体
		switch field.Type.Kind() {
		case reflect.Struct:
			// 如果是匿名嵌入的结构体（匿名字段），递归处理其字段
			if field.Anonymous {
				embeddedValue := fieldValue
				if embeddedValue.Kind() == reflect.Ptr {
					if embeddedValue.IsNil() {
						continue // 跳过 nil 指针
					}
					embeddedValue = embeddedValue.Elem()
				}
				buildFieldMapping(embeddedValue, field.Type, prefix, tagToField, fieldValueMap)
			} else {
				// 如果是非匿名结构体字段，将其作为一个整体字段处理
				(*tagToField)[columnName] = field
				(*fieldValueMap)[columnName] = fieldValue
			}
		case reflect.Ptr:
			if field.Type.Elem().Kind() == reflect.Struct {
				// 指向结构体的指针
				if !fieldValue.IsNil() {
					embeddedValue := fieldValue.Elem()
					if field.Anonymous { // 匿名嵌入的指针结构体
						buildFieldMapping(embeddedValue, field.Type.Elem(), prefix, tagToField, fieldValueMap)
					} else { // 非匿名指针结构体字段
						(*tagToField)[columnName] = field
						(*fieldValueMap)[columnName] = fieldValue
					}
				} else if field.Anonymous {
					// 匿名指针字段为 nil，跳过
					continue
				} else {
					// 非匿名指针字段为 nil，作为普通字段处理
					(*tagToField)[columnName] = field
					(*fieldValueMap)[columnName] = fieldValue
				}
			} else {
				// 非结构体指针字段，作为普通字段处理
				(*tagToField)[columnName] = field
				(*fieldValueMap)[columnName] = fieldValue
			}
		default:
			// 普通字段
			(*tagToField)[columnName] = field
			(*fieldValueMap)[columnName] = fieldValue
		}
	}
}

// 从 gorm tag 中提取 column 名称
func getColumnFromGormTag(gormTag string) string {
	if gormTag == "" {
		return ""
	}

	// 按分号分割 tag 选项
	options := strings.Split(gormTag, ";")
	for _, option := range options {
		option = strings.TrimSpace(option)
		// 查找 column:xxx 格式的选项
		if strings.HasPrefix(option, "column:") {
			return strings.TrimPrefix(option, "column:")
		}
	}
	return ""
}

type FilterParam struct {
	query any
	args  []any
}

func Filter(query any, args ...any) *FilterParam {
	return &FilterParam{query: query, args: args}
}

// UpdateWithMap 根据过滤条件和map参数更新
// ctx.DB().UpdateWithMap(&model.BlueUserSaas{}, app.Filter("id = ?", id), map[string]any{"status": status})
func (d *Database) UpdateWithMap(model Moduler, f *FilterParam, data map[string]any) Result {

	res := d.db().Model(model).Where(f.query, f.args...).Updates(data)

	var updateInfos []string
	for k, v := range data {
		updateInfos = append(updateInfos, fmt.Sprintf("%s=%v", k, v))
	}
	d.Log().InfofSkip(1, "更新表:%s, 数据:%s", model.TableName(), strings.Join(updateInfos, ","))
	return Result{RowsAffected: res.RowsAffected, Error: rerr.Wrap(res.Error)}
}
