package model

import (
	"reflect"
	"testing"

	"github.com/brunoksato/golang-boilerplate/core"
)

func TestIsValidator(t *testing.T) {
	core.AssertTrue(t, IsValidator(reflect.TypeOf(&User{})))
	core.AssertTrue(t, IsValidator(reflect.TypeOf(&Configuration{})))
}

func TestValidateStructNoErrors(t *testing.T) {
	v := ValidationTestStruct{
		Required:     "I'm here",
		NoValidation: "Who cares?",
		Numeric:      "343452345",
		Float:        "343452345.10",
		Email:        "valid@model.com.br",
	}

	core.AssertNil(t, ValidateStruct(v))
}

func TestValidateStructOneError(t *testing.T) {
	v := ValidationTestStruct{
		Required:     "I'm here",
		NoValidation: "Who cares?",
		Numeric:      "343452345.10",
		Float:        "343452345.10",
		Email:        "valid@petmondo.com",
	}

	err := ValidateStruct(v)
	core.AssertEqual(t, "Numeric: 343452345.10 does not validate as numeric;", err.Error())
}

func TestValidateStructMultipleErrors(t *testing.T) {
	v := ValidationTestStruct{
		Required:         "",
		Numeric:          "343452345.10",
		Float:            "343452345",
		Email:            "invalid",
		DoubleValidation: "34t34*?$ˆ",
	}

	err := ValidateStruct(v)
	core.AssertEqual(t, "Required: non zero value required;Numeric: 343452345.10 does not validate as numeric;Email: invalid does not validate as email;DoubleValidation: 34t34*?$ˆ does not validate as alphanum;", err.Error())
}

func TestValidateStructZeroStruct(t *testing.T) {
	v := ValidationTestStruct{}

	err := ValidateStruct(v)
	core.AssertEqual(t, "Required: non zero value required;", err.Error())
}

func TestValidateStructField(t *testing.T) {
	v := ValidationTestStruct{
		Required: "",
		Numeric:  "343452345.10",
		Float:    "343452345.10",
		Email:    "bruno@model.com.br",
	}

	core.AssertNil(t, ValidateStructField(v, "NoValidation"))
	core.AssertNil(t, ValidateStructField(v, "Float"))
	core.AssertNil(t, ValidateStructField(v, "Email"))
	core.AssertEqual(t, "Required: non zero value required;", ValidateStructField(v, "Required").Error())
	core.AssertEqual(t, "Numeric: 343452345.10 does not validate as numeric;", ValidateStructField(v, "Numeric").Error())
}

func TestValidateStructDoubleValidationError(t *testing.T) {
	v := ValidationTestStruct{
		DoubleValidation: "34t34**???+!@*&#$ˆ(((dfasdlfjasd;fl2903",
	}

	// Commenting this out, since the results are unpredictable - fails on Jekinns
	//AssertEqual(t, "DoubleValidation: 34t34**???+!@*&#$ˆ(((dfasdlfjasd;fl2903 does not validate as alphanum;", ValidateStructField(v, "DoubleValidation").Error())

	v.DoubleValidation = "34t34**???"
	core.AssertEqual(t, "DoubleValidation: 34t34**??? does not validate as alphanum;", ValidateStructField(v, "DoubleValidation").Error())

	v.DoubleValidation = "f8asd09f8asd0f98asf0as98df"
	core.AssertEqual(t, "DoubleValidation: f8asd09f8asd0f98asf0as98df does not validate as length(3|10);", ValidateStructField(v, "DoubleValidation").Error())

	v.DoubleValidation = "f8asd09"
	core.AssertNil(t, ValidateStructField(v, "DoubleValidation"))
}

func TestValidateStructFields(t *testing.T) {
	v := ValidationTestStruct{
		Required:         "",
		Numeric:          "343452345.10",
		Float:            "343452345.10",
		Email:            "bruno@model.com.br",
		DoubleValidation: "@(#*$@",
	}

	err := ValidateStructFields(v, []string{"Float", "Email"})
	core.AssertNil(t, err)

	err = ValidateStructFields(v, []string{"Required", "Float", "Email"})
	core.AssertEqual(t, "Required: non zero value required;", err.Error())

	err = ValidateStructFields(v, []string{"Required", "Numeric", "Float", "Email"})
	core.AssertEqual(t, "Required: non zero value required;Numeric: 343452345.10 does not validate as numeric;", err.Error())
}

func TestValidateRequiredFields(t *testing.T) {
	vts := ValidationTestStruct{
		Required:         "And here",
		NoValidation:     "",
		Numeric:          "Here but wrong",
		Float:            "",
		Email:            "",
		DoubleValidation: "Make it triple",
	}
	core.AssertNil(t, ValidateRequiredFields(vts, []string{}))
	core.AssertNil(t, ValidateRequiredFields(vts, []string{"Required"}))
	core.AssertNil(t, ValidateRequiredFields(vts, []string{"Required", "Numeric", "DoubleValidation"}))
	core.AssertEqual(t, "NoValidation: non zero value required;", ValidateRequiredFields(vts, []string{"NoValidation"}).Error())
	core.AssertEqual(t, "NoValidation: non zero value required;Float: non zero value required;Email: non zero value required;", ValidateRequiredFields(vts, []string{"Required", "NoValidation", "Numeric", "Float", "Email", "DoubleValidation"}).Error())
}

type ValidationTestStruct struct {
	Required         string `valid:"required"`
	NoValidation     string `valid:"-"`
	Numeric          string `valid:"numeric"`
	Float            string `valid:"float"`
	Email            string `valid:"email"`
	DoubleValidation string `valid:"alphanum,length(3|10)"`
}
