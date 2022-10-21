package util

func UnwrapPtr[IT any, OT any](fn func(input *IT) *OT) func(input *IT, err error) (*OT, error) {
	return func(input *IT, err error) (*OT, error) {
		if err != nil {
			return nil, err
		}
		if input == nil {
			return nil, nil
		}
		return fn(input), nil
	}
}
