package model

type Model[T any] struct {
	model T
}

func GetModelStructUsingGenerics[T any](tableName string) *Model[T] {
	var modelStruct interface{}

	switch tableName {
	case "video_release":
		modelStruct = VideoRelease{}
	}

	return &Model[T]{model: modelStruct.(T)}
}

func GetModelStruct(tableName string) interface{} {
	var modelStruct interface{}

	switch tableName {
	case "video_release":
		modelStruct = VideoRelease{}
	}

	return modelStruct
}
