package model

import (
	"reflect"
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
	"github.com/jinzhu/gorm"
)

func TestSorterOverriding(t *testing.T) {
	setupDB()
	defer teardownDB()
	setupSorterDB()

	s1 := SorterTestStruct{Order: 1}
	s2 := SorterTestStruct{Order: 2}
	TESTDB.Create(&s1)
	TESTDB.Create(&s2)

	var sorters []SorterTestStruct
	db := s1.OrderBy(TESTDB)
	db.Find(&sorters)
	core.AssertEqual(t, 2, len(sorters))
	core.AssertEqual(t, s1.ID, sorters[0].ID)
	core.AssertEqual(t, s2.ID, sorters[1].ID)
}

func TestIsSorter(t *testing.T) {
	core.AssertTrue(t, IsSorter(reflect.TypeOf(User{})))
	core.AssertTrue(t, IsSorter(reflect.TypeOf(Configuration{})))
}

type SorterTestStruct struct {
	Model
	Order uint
}

func (m SorterTestStruct) OrderBy(db *gorm.DB) *gorm.DB {
	order := "\"order\" asc"
	return db.Order(order)
}

func setupSorterDB() {
	TESTDB.AutoMigrate(&SorterTestStruct{})
}
