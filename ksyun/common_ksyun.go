package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"log"
	"reflect"
	"strings"
)

//Convert sdk response type (map[string]interface{}) to the type terraform can realized([]map[string]interface).
//params data limit ： [k,v]:the type of k must be string ,the type of v must be basic type.
func GetSubDByRep(data interface{}, include, exclude map[string]bool) []interface{} {
	ma, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	subD := make(map[string]interface{})
	for k, v := range ma {
		if exclude[k] || !include[k] {
			continue
		}
		subD[Hump2Downline(k)] = v
	}
	return []interface{}{subD}
}

//sdk resp []map[string]interface{}->terraform schema.ResourceData
//Convert sdk response type ([]map[string]interface{}) to the type terraform can realized([]map[string]interface).
//include ：representing the key terraform has defined.
//exclude ：representing the key which the type is not basic type.
//Suitable for the value in d.Set（ key，value）,and the type of value must be List.
func GetSubSliceDByRep(items []interface{}, include /*,exclude*/ map[string]bool) []map[string]interface{} {
	datas := []map[string]interface{}{}
	for _, v := range items {
		data := map[string]interface{}{}
		vv, _ := v.(map[string]interface{})
		for key, value := range vv {
			//ignore keys whose type is not basic type,and need to deal later.
			if /*exclude[key]||*/ !include[key] {
				continue //if not judge,sdk may set value to terraform which can identify,and will panic.
			}
			data[Hump2Downline(key)] = value
		}
		datas = append(datas, data)
	}
	return datas
}

//sdk resp map[string]interface{} inline struct ->terraform schema.ResourceData
//convert inline struct from sdk response type ([]map[string]interface{}) to the type terraform can realized([]map[string]interface).
//exclude ：representing the key which the type is not basic type.
func GetSubStructDByRep(datas interface{}, exclude map[string]bool) map[string]interface{} {

	subStruct := map[string]interface{}{}
	items, ok := datas.(map[string]interface{})
	if !ok {
		return subStruct
	}
	for k, v := range items {
		if exclude[k] {
			continue
		}
		subStruct[Hump2Downline(k)] = v
	}
	return subStruct
}

//set sdk response (map[string]interface{}) to the terraform ([]map[string]interface).
//params data limit ： [k,v]:the type of k must be string ,the type of v must be basic type.
//exclude ：representing the key which the type is not basic type (terraform can't identity the type which is not basic type).
//mre: the params not set to terraform .
func SetDByRespV1(d *schema.ResourceData, m interface{}, exclud map[string]bool) map[string]interface{} {
	ma, ok := m.(map[string]interface{})
	mre := make(map[string]interface{})
	if !ok {
		return mre
	}
	for k, v := range ma {
		if exclud[k] {
			if mm, ok := v.(map[string]interface{}); ok {
				mre[k] = mm
			} else {
				mre[k] = v
			}
			continue
		}
		err := d.Set(Hump2Downline(k), v)
		if err != nil {
			log.Println("ERROR: SetDByRespV1 failed:", err.Error())
			panic("ERROR: SetDByRespV1 failed:" + err.Error())
			//return mre
		}
	}
	return mre
}

func SetResourceDataByResp(d *schema.ResourceData, item interface{}, keys map[string]SdkResponseData, start ...bool) interface{} {
	setFlag := false
	if start == nil || (len(start) > 0 && start[0]) {
		setFlag = true
	}
	if reflect.ValueOf(item).Kind() == reflect.Map {
		result := make(map[string]interface{})
		root := item.(map[string]interface{})
		for k, v := range root {
			var value interface{}
			var err error
			if keys[k].Next != nil {
				value = SetResourceDataByResp(d, v, keys[k].Next, false)
			} else if _, ok := keys[k]; ok {
				value = v
			} else {
				continue
			}
			if setFlag {
				err = d.Set(Hump2Downline(k), value)
				if err != nil {
					log.Println(err.Error())
					panic("ERROR: " + err.Error())
				}
			} else {
				result[Hump2Downline(k)] = value
			}
		}
		if len(result) > 0 {
			return result
		}
	} else if reflect.ValueOf(item).Kind() == reflect.Slice {
		var result []interface{}
		result = []interface{}{}
		root := item.([]interface{})
		for _, v := range root {
			value := SetResourceDataByResp(d, v, keys, false)
			result = append(result, value)
		}
		if len(result) > 0 {
			return result
		}
	}
	return nil
}

func addProjectInfoAll(d *schema.ResourceData, req *map[string]interface{}, client *KsyunClient, key ...string) error {
	var project string
	if key != nil && len(key) == 1 {
		project = key[0]
	} else {
		project = "project_id"
	}
	return getProjectInfo(req, client, project)
}

func addProjectInfo(d *schema.ResourceData, req *map[string]interface{}, client *KsyunClient, key ...string) error {
	var project string
	if key != nil && len(key) == 1 {
		project = key[0]
	} else {
		project = "project_id"
	}
	hump := Downline2Hump(project)
	if pj, ok := d.GetOk(project); ok {
		(*req)[hump+".1"] = fmt.Sprintf("%v", pj)
	} else if d.HasChange(project) && !d.IsNewResource() {
		(*req)[hump+".1"] = fmt.Sprintf("%v", d.Get(project))
	} else {
		return getProjectInfo(req, client, project)
	}
	return nil
}

func checkStringValueNotEmpty(d *schema.ResourceData, key string) (string, bool) {
	if v, ok := d.GetOk(key); ok && v != nil && v.(string) != "" {
		return v.(string), true
	}
	return "", false
}

func getProjectInfo(input *map[string]interface{}, client *KsyunClient, key ...string) error {
	var project string
	if key != nil && len(key) == 1 {
		project = key[0]
	} else {
		project = "project_id"
	}
	hump := Downline2Hump(project)
	iamConn := client.iamconn
	req := make(map[string]interface{})
	resp, err := iamConn.GetAccountAllProjectList(&req)
	if err != nil {
		return fmt.Errorf("Error GetAccountAllProjectList : %s ", err)
	}
	if resp != nil {
		l, err1 := getSdkValue("ListProjectResult.ProjectList", *resp)
		if err1 != nil {
			return fmt.Errorf("Error GetAccountAllProjectList : %s ", err)
		}
		if l != nil {
			if l1, ok := l.([]interface{}); ok {
				for i, pj := range l1 {
					index := fmt.Sprintf("%v.%d", hump, i+1)
					(*input)[index] = pj.(map[string]interface{})["ProjectId"]
				}
			}
		}
	}
	return nil
}

//set sdk response (map[string]interface{}) to the terraform ([]map[string]interface).
//params data limit ： [k,v]:the type of k must be string ,the type of v must be basic type.
//include ：representing the key terraform has defined. terraform will panic if set the key that not defined.
//exclude ：representing the key which the type is not basic type (terraform can't identity the type which is not basic type).
//mre: the params not set to terraform .
func SetDByResp(d *schema.ResourceData, m interface{}, includ, exclude map[string]bool) map[string]interface{} {
	mre := make(map[string]interface{})
	ma, ok := m.(map[string]interface{})
	if !ok {
		return mre
	}
	for k, v := range ma {
		if !includ[k] || exclude[k] {
			if mm, ok := v.(map[string]interface{}); ok {
				mre[k] = mm
			} else {
				mre[k] = v
			}
			continue
		}

		err := d.Set(Hump2Downline(k), v)
		if err != nil {
			log.Println(err.Error())
			panic("ERROR: " + err.Error())
		}
	}
	return mre
}

//The hump is converted to an underline simply, and no special treatment is required for even uppercase letters.
//ex:aDDCC ->a_d_d_c_c
func Hump2Downline(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	var s1 string
	if len(s) == 1 {
		s1 = strings.ToLower(s[:1])
		return s1
	}
	for k, v := range s {
		if k == 0 {
			s1 = strings.ToLower(s[0:1])
			continue
		}
		if v >= 65 && v <= 90 {
			v1 := "_" + strings.ToLower(s[k:k+1])
			s1 = s1 + v1
		} else {
			s1 = s1 + s[k:k+1]
		}
	}
	return s1
}

//The underline is converted to an hump simply.
func Downline2Hump(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	var s1 []string
	ss := strings.Split(s, "_")
	for _, v := range ss {
		vv := strings.ToUpper(v[:1]) + v[1:]
		s1 = append(s1, vv)
	}
	return strings.Join(s1, "")
}

func Downline2Filter(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	var s1 []string
	ss := strings.Split(s, "_")
	for i, v := range ss {
		if i < len(ss)-1 {
			vv := v + "-"
			s1 = append(s1, vv)
		} else {
			s1 = append(s1, v)
		}
	}
	return strings.Join(s1, "")
}

//flattern struct
// convert input param struct to map when create(with out prefix).
func FlatternStruct(v interface{}, req *map[string]interface{}) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				if len(fmt.Sprintf("%v", v2)) == 0 {
					continue
				}
				vv := Downline2Hump(k2)
				(*req)[vv] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//flattern struct Suitable for inline struct
//convert input param struct to map when create(with  prefix).
//prefix: the name of the outer structure
func FlatternStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				if len(fmt.Sprintf("%v", v2)) == 0 {
					continue
				}
				kk := Downline2Hump(k2)
				(*req)[fmt.Sprintf("%s.%s", prex, kk)] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//FlatternStructSlicePrefix 用于创建时，结构体切片类型的入参转换为map型 ,【
//Flattern StructSlice Suitable for the slice of inline struct
//convert input param struct to map when create(with  prefix).
//prefix: the name of the slice
func FlatternStructSlicePrefix(values interface{}, req *map[string]interface{}, prex string) {
	v, _ := values.([]interface{})
	k := 0
	for _, v1 := range v {
		vv := v1.(map[string]interface{})
		if len(vv) == 0 {
			continue
		}
		k++
		for k2, v2 := range vv {
			kk := Downline2Hump(k2)
			(*req)[fmt.Sprintf("%s.%d.%s", prex, k, kk)] = fmt.Sprintf("%v", v2)
		}
	}
}

//Suitable for filter which need conver param with "_"(terraform) to "-"(sdk) when read .
//convert input param struct to map when create(without prefix).
func ConvertFilterStruct(v interface{}, req *map[string]interface{}) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				if len(fmt.Sprintf("%v", v2)) == 0 {
					continue
				}
				vv := strings.ReplaceAll(k2, "_", "-")
				(*req)[vv] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//Suitable for filter which need conver param with "_"(terraform) to "-"(sdk) when read.
//convert input param struct to map when create(with prefix).
//prefix:the name of the elemet from filter
func ConvertFilterStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			if v1[0] == nil {
				return
			}
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				if len(fmt.Sprintf("%v", v2)) == 0 {
					continue
				}
				vv := strings.ReplaceAll(k2, "_", "-")
				(*req)[fmt.Sprintf("%s.%s", prex, vv)] = v2
			}
		}
	}
}

/*
func ConvertFilterStructStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			if v1[0] == nil {
				return
			}
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				vv := strings.ReplaceAll(k2, "-", "_")
				v3, ok3 := v2.([]string)
				if !ok3 || len(v3) == 0 {
					(*req)[fmt.Sprintf("%s.%s", prex, vv)] = fmt.Sprintf("%v", v2)
				}
				(*req)[fmt.Sprintf("%s.%s", prex, vv)] = fmt.Sprintf("%v", v3[0])

			}
		}
	}
}

*/
func dataSourceKscSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {

	d.SetId(hashStringArray(ids))
	if err := d.Set("total_count", len(datas)); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if err := d.Set(dataKey, datas); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		return writeToFile(outputFile.(string), datas)
	}

	return nil
}
func dataSourceKscSaveSlice(d *schema.ResourceData, dataKey string, ids []string, datas []string) error {

	d.SetId(hashStringArray(ids))
	if err := d.Set("total_count", len(datas)); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}

	if err := d.Set(dataKey, datas); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		return writeToFile(outputFile.(string), datas)
	}

	return nil
}

func dataSourceDbSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {
	if len(ids) == 1 {
		d.SetId(ids[0])
	} else {
		d.SetId(strings.Join(ids, ","))
	}

	//if err := d.Set("total_count", len(datas)); err != nil {
	//	return fmt.Errorf("error set datas %v :%v", datas, err)
	//}
	log.Printf("reset  dataKey: %v datas: %v", dataKey, datas)

	if err := d.Set(dataKey, datas); err != nil {
		logger.DebugInfo("err %+v", err)
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		logger.DebugInfo(" output file name : %+v", outputFile.(string)+"_"+d.Id())
		return writeToFile(outputFile.(string)+"_"+d.Id(), datas)
	} else {
		return fmt.Errorf(" output file error,  %+v", outputFile)
	}
}

func dataDbSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {
	if len(ids) == 1 {
		d.SetId(ids[0])
	} else {
		d.SetId(strings.Join(ids, ","))
	}

	if err := d.Set("total_count", len(datas)); err != nil {
		return err
	}
	log.Printf("reset  dataKey: %v datas: %v", dataKey, datas)

	if err := d.Set(dataKey, datas); err != nil {
		logger.DebugInfo("$$err$$ %+v", err)
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		logger.DebugInfo(" ------------ %+v", outputFile)
		return writeToFile(outputFile.(string)+"_"+"data", datas)
	} else {
		return fmt.Errorf(" !!! %+v", outputFile)
	}
}

func SetDByFkResp(d *schema.ResourceData, m interface{}, include map[string]bool) map[string]interface{} {
	exclude := make(map[string]interface{})
	ma, ok := m.(map[string]interface{})
	if !ok {
		return exclude
	}
	for k, v := range ma {
		if !include[k] {
			if mm, ok := v.(map[string]interface{}); ok {
				exclude[k] = mm
			} else {
				exclude[k] = v
			}
			continue
		}

		err := d.Set(k, v)
		if err != nil {
			log.Println(err.Error())
			panic("ERROR: " + err.Error())
		}
	}
	logger.DebugInfo("-----exclude fields : %+v", exclude)
	return exclude
}

func Camel2Hungarian(s string) string {
	in := true
	var s1 string
	for k, v := range s {
		if v >= 65 && v <= 90 {
			if !in || (k < len(s)-1 && k > 0 && s[k+1] >= 97 && s[k+1] <= 122) {
				s1 += "_" + string(v+32)
				in = true
			} else {
				s1 += string(v + 32)
			}
		} else if v == 46 {
			if k < len(s)-1 && s[k+1] >= 48 && s[k+1] <= 57 {
				s1 += "_"
				in = false
			} else {
				s1 += string(v)
				in = true
			}
		} else {
			s1 += string(v)
			in = false
		}
	}

	return s1
}
