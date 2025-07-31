package ksyun

type (
	pageCall              func(map[string]interface{}) ([]interface{}, error)
	pageCallWithNextToken func(map[string]interface{}) ([]interface{}, string, error)
)

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

func pageQueryWithNextToken(condition map[string]interface{}, limitParam string, nextTokenParam string, limit int, call pageCallWithNextToken) (data []interface{}, err error) {
	if condition == nil {
		condition = make(map[string]interface{})
	}
	nextToken := ""
	for {
		var d []interface{}
		condition[limitParam] = limit
		if nextToken != "" {
			condition[nextTokenParam] = nextToken
		}
		d, nextToken, err = call(condition)
		if err != nil {
			return data, err
		}
		data = append(data, d...)
		if len(d) < limit {
			break
		}
	}
	return data, err
}
