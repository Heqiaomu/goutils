package conform

// Combine 将source作为target的组成部分拼接上去
func Combine(source *UIData, target *UIData) *UIData {
	target.Fields = append(target.Fields, source.Fields...)
	target.Inputs = append(target.Inputs, source.Inputs...)
	target.FieldReactions = append(target.FieldReactions, source.FieldReactions...)
	id2UIDataTarget := target.Id2UIData
	id2UIDataSource := source.Id2UIData
	for k, v := range id2UIDataSource.InputKey2Input {
		id2UIDataTarget.InputKey2Input[k] = v
	}
	for k, v := range id2UIDataSource.FieldKey2Field {
		id2UIDataTarget.FieldKey2Field[k] = v
	}
	for k, v := range id2UIDataSource.FieldKey2Reactions {
		id2UIDataTarget.FieldKey2Reactions[k] = v
	}
	return target
}
