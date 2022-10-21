package util

func UnwrapPtr[IT any, OT any](fn func(input *IT) OT) func(input *IT, err error) (OT, error) {
	var zv OT
	return func(input *IT, err error) (OT, error) {
		if err != nil {
			return zv, err
		}
		if input == nil {
			return zv, nil
		}
		return fn(input), nil
	}
}
