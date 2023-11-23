package ksyun

type pageCall func(map[string]interface{}) ([]interface{}, error)

func pageQuery(condition map[string]interface{}, limitParam string, pageParam string, limit int, start int, call pageCall) (data []interface{}, err error) {
	if condition == nil {
		condition = make(map[string]interface{})
	}
	offset := start
	for {
		var d []interface{}
		condition[limitParam] = limit
		condition[pageParam] = offset
		d, err = call(condition)
		if err != nil {
			return data, err
		}
		data = append(data, d...)
		if len(d) < limit {
			break
		}
		offset = offset + limit
	}
	return data, err
}
