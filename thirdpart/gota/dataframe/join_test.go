package dataframe

import (
	"fmt"
	"testing"

	df "github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/stretchr/testify/assert"
)

// 构建两个 `DataFrame` 实例
//
// 其中 `employee` 实例 (左) 的 `Department` 字段和 `department` 实例 (右) 的 `Id` 字段用于连接
func makeDFForJoin(t *testing.T) (df.DataFrame, df.DataFrame) {
	employee := LoadSliceMap(map[string][]any{
		"Id":         {1, 2, 3, 4, 5, 6},
		"Name":       {"Alvin", "Emma", "Lucy", "Tom", "Arthur", "Jenny"},
		"Gender":     {"M", "F", "F", "M", "M", "F"},
		"Age":        {42, 38, 51, 34, 21, 33},
		"Salaries":   {32000, 18000, 9500, 10000, 7850, 9200},
		"Department": {1, 2, 1, 2, 1, 3},
	})
	assert.NoError(t, employee.Error())

	department := LoadSliceMap(map[string][]any{
		"Id":   {1, 2, 4},
		"Name": {"R&D", "Sales", "HR"},
	})
	assert.NoError(t, department.Error())

	fmt.Println(employee, department)
	return employee, department
}

// 测试两个 `DataFrame` 的内连接操作
//
// 内连接即 `InnerJoin`, 结果为一个包含两个 `DataFrame` 具备相同指定字段值的记录行的 `DataFrame`
//
// 通过 `Department_Id` 字段 (在 `employee` 实例中为 `Department` 字段, 在 `department`
// 实例中为 `Id` 字段) 进行连接连接结果通过 `Employee_Name` 字段值为 `"Alvin"` 进行过滤
func TestDataFrameJoin_InnerJoin(t *testing.T) {
	employee, department := makeDFForJoin(t)

	rs := employee.
		Rename("Employee_Id", "Id").
		Rename("Department_Id", "Department").
		Rename("Employee_Name", "Name").
		InnerJoin(
			department.
				Rename("Department_Id", "Id").
				Rename("Department_Name", "Name"),
			"Department_Id",
		).
		Filter(df.F{Colname: "Employee_Name", Comparator: series.Eq, Comparando: "Alvin"})

	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 查看过滤结果, `Employee_Name` 字段值为 "Alvin" 的 `Department_Name` 为 `"R&D"`
	assert.Equal(t, 1, rs.Nrow())
	assert.Equal(t, "R&D", rs.Col("Department_Name").Val(0))
}

// 测试两个 `DataFrame` 的交叉连接操作
//
// 交叉连接即 `CrossJoin`, 结果为一个包含两个 `DataFrame` 具备相同指定字段值的记录行的 `DataFrame`
//
// 通过 `Department_Id` 字段 (在 `employee` 实例中为 `Department` 字段, 在 `department`
// 实例中为 `Id` 字段) 进行连接, 连接结果通过 `Employee_Name` 字段值为 `"Alvin"` 进行过滤
func TestDataFrameJoin_CrossJoin(t *testing.T) {
	employee, department := makeDFForJoin(t)

	rs := employee.
		Rename("Employee_Id", "Id").
		Rename("Department_Id", "Department").
		Rename("Employee_Name", "Name").
		CrossJoin(
			department.
				Rename("Department_Id", "Id").
				Rename("Department_Name", "Name"),
		).
		Filter(df.F{Colname: "Employee_Name", Comparator: series.Eq, Comparando: "Alvin"})

	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 查看过滤结果, `Employee_Name` 字段值为 "Alvin" 的 `Department_Name` 字段值为 `"R&D"`, `"Sales"` 和 `"HR"`
	// 即 `employee` 实例的所有行和 `depart` 实例的所有行交叉连接
	assert.Equal(t, 3, rs.Nrow())
	assert.Equal(t, []string{"R&D", "Sales", "HR"}, rs.Col("Department_Name").Records())
}

// 测试两个 `DataFrame` 的外连接操作
//
// 外连接即 `OuterJoin`, 结果为一个包含两个 `DataFrame` 具备相同指定字段值的记录行的 `DataFrame`
//
// 通过 `Department_Id` 字段 (在 `employee` 实例中为 `Department` 字段, 在 `department`
// 实例中为 `Id` 字段) 进行连接连接结果通过 `Employee_Name` 字段值为 `"Alvin"` 或者 `Department_Name`
// 字段值为 `"Jenny"` 进行过滤
func TestDataFrameJoin_OuterJoin(t *testing.T) {
	employee, department := makeDFForJoin(t)

	rs := employee.
		Rename("Employee_Id", "Id").
		Rename("Department_Id", "Department").
		Rename("Employee_Name", "Name").
		OuterJoin(
			department.
				Rename("Department_Id", "Id").
				Rename("Department_Name", "Name"),
			"Department_Id",
		).
		Filter(
			df.F{Colname: "Department_Name", Comparator: series.Eq, Comparando: "HR"},
			df.F{Colname: "Employee_Name", Comparator: series.Eq, Comparando: "Jenny"},
		)

	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 查看过滤结果:
	//  - `Employee_Name` 字段值为 "Jenny" 的 `Department_Name` 字段值为 `nil`
	//  - `Department_Name` 字段值为 "HR" 的 `Employee_Name` 字段值为 `nil`
	// 即 `employee` 实例中 `Employee_Name` 为 `Jenny` 的行在 `depart` 实例中不存在对应行, `depart` 实例中 `Department_Name` 为 `HR` 的行在 `employee` 实例中不存在对应行
	assert.Equal(t, 2, rs.Nrow())
	assert.Equal(t, []string{"NaN", "HR"}, rs.Col("Department_Name").Records())
	assert.Equal(t, []string{"Jenny", "NaN"}, rs.Col("Employee_Name").Records())
}

// 测试两个 `DataFrame` 的左外连接操作
//
// 左外连接即 `LeftJoin`, 结果为一个包含两个 `DataFrame` 具备相同指定字段值的记录行的 `DataFrame`
//
// 通过 `Department_Id` 字段 (在 `employee` 实例中为 `Department` 字段, 在 `department`
// 实例中为 `Id` 字段) 进行连接连接结果通过 `Employee_Name` 字段值为 `"Alvin"` 或者 `Department_Name`
// 字段值为 `"Jenny"` 进行过滤
func TestDataFrameJoin_LeftJoin(t *testing.T) {
	employee, department := makeDFForJoin(t)

	rs := employee.
		Rename("Employee_Id", "Id").
		Rename("Department_Id", "Department").
		Rename("Employee_Name", "Name").
		LeftJoin(
			department.
				Rename("Department_Id", "Id").
				Rename("Department_Name", "Name"),
			"Department_Id",
		).
		Filter(
			df.F{Colname: "Department_Name", Comparator: series.Eq, Comparando: "HR"},
			df.F{Colname: "Employee_Name", Comparator: series.Eq, Comparando: "Jenny"},
		)

	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 查看过滤结果:
	//  - `Employee_Name` 字段值为 "Jenny" 的 `Department_Name` 字段值为 `nil`
	//  - `Department_Name` 字段值为 "HR" 的 `Employee_Name` 字段值为 `nil`
	// 即 `employee` 实例中 `Employee_Name` 为 `Jenny` 的行在 `depart` 实例中不存在对应行, `depart` 实例中 `Department_Name` 为 `HR` 的行在 `employee` 实例中不存在对应行
	assert.Equal(t, 1, rs.Nrow())
	assert.Equal(t, []string{"NaN"}, rs.Col("Department_Name").Records())
	assert.Equal(t, []string{"Jenny"}, rs.Col("Employee_Name").Records())
}

// 测试两个 `DataFrame` 的右外连接操作
//
// 右外连接即 `RightJoin`, 结果为一个包含两个 `DataFrame` 具备相同指定字段值的记录行的 `DataFrame`
//
// 通过 `Department_Id` 字段 (在 `employee` 实例中为 `Department` 字段, 在 `department`
// 实例中为 `Id` 字段) 进行连接连接结果通过 `Employee_Name` 字段值为 `"Alvin"` 或者 `Department_Name`
// 字段值为 `"Jenny"` 进行过滤
func TestDataFrameJoin_RightJoin(t *testing.T) {
	employee, department := makeDFForJoin(t)

	rs := employee.
		Rename("Employee_Id", "Id").
		Rename("Department_Id", "Department").
		Rename("Employee_Name", "Name").
		RightJoin(
			department.
				Rename("Department_Id", "Id").
				Rename("Department_Name", "Name"),
			"Department_Id",
		).
		Filter(
			df.F{Colname: "Department_Name", Comparator: series.Eq, Comparando: "HR"},
			df.F{Colname: "Employee_Name", Comparator: series.Eq, Comparando: "Jenny"},
		)

	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 查看过滤结果:
	//  - `Employee_Name` 字段值为 "Jenny" 的 `Department_Name` 字段值为 `nil`
	//  - `Department_Name` 字段值为 "HR" 的 `Employee_Name` 字段值为 `nil`
	// 即 `employee` 实例中 `Employee_Name` 为 `Jenny` 的行在 `depart` 实例中不存在对应行, `depart` 实例中 `Department_Name` 为 `HR` 的行在 `employee` 实例中不存在对应行
	assert.Equal(t, 1, rs.Nrow())
	assert.Equal(t, []string{"HR"}, rs.Col("Department_Name").Records())
	assert.Equal(t, []string{"NaN"}, rs.Col("Employee_Name").Records())
}
