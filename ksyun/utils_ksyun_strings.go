package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func stringSplitDiffSuppressFunc(sep string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		olds := strings.Split(old, sep)
		news := strings.Split(new, sep)
		count := 0
		if len(olds) != len(news) {
			return false
		}
		for _, o := range olds {
			for _, n := range news {
				if o == n {
					count = count + 1
				}
			}
		}
		if count == len(olds) {
			return true
		}
		return false
	}
}

func stringSplitSchemaValidateFunc(sep string) schema.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}
		for _, s := range strings.Split(v, sep) {
			if strings.Count(v, s) > 1 {
				errors = append(errors, fmt.Errorf(" %s only allow one in %s", s, k))
			}
		}
		return warnings, errors
	}
}

func stringSplitChange(sep string, field string, currentExists []string, d *schema.ResourceData) (add []string, del []string) {
	var (
		oldVal interface{}
		newVal interface{}
	)
	if d.HasChange(field) {
		oldVal, newVal = d.GetChange(field)
		for _, current := range currentExists {
			exist := false
			for _, change := range strings.Split(newVal.(string), sep) {
				if change == current {
					exist = true
					break
				}
			}
			if !exist && strings.Contains(oldVal.(string), current) {
				del = append(del, current)
			}
		}
		for _, change := range strings.Split(newVal.(string), sep) {
			exist := false
			for _, current := range currentExists {
				if current == change {
					exist = true
					break
				}
			}
			if !exist {
				add = append(add, change)
			}
		}
	}
	return add, del
}

func stringSplitRead(sep string, field string, currentExists []string, d *schema.ResourceData) (result string) {
	for _, current := range currentExists {
		result = result + current + sep
	}
	if result != "" {
		result = result[0 : len(result)-1]
	}
	return result
}

func stringSliceContains(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}
