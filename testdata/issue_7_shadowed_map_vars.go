package testdata

type SomeStruct struct {
	mapSlice map[string][]string
}

type SomeStruct2 struct {
	mapStruct map[string]SomeStruct
}

// DeepCopy generates a deep copy of SomeStruct2
func (o SomeStruct2) DeepCopy() SomeStruct2 {
	var cp SomeStruct2 = o
	if o.mapStruct != nil {
		cp.mapStruct = make(map[string]SomeStruct, len(o.mapStruct))
		for k, v := range o.mapStruct {
			var cp_mapStruct_v SomeStruct
			if v.mapSlice != nil {
				cp_mapStruct_v.mapSlice = make(map[string][]string, len(v.mapSlice))
				for k, v := range v.mapSlice {
					var cp_mapStruct_v_mapSlice_v []string
					if v != nil {
						cp_mapStruct_v_mapSlice_v = make([]string, len(v))
						copy(cp_mapStruct_v_mapSlice_v, v)
					}
					cp_mapStruct_v.mapSlice[k] = cp_mapStruct_v_mapSlice_v
				}
			}
			cp.mapStruct[k] = cp_mapStruct_v
		}
	}
	return cp
}
