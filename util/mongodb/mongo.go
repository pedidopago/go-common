package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

func E(key string, value interface{}) primitive.E {
	return primitive.E{
		Key:   key,
		Value: value,
	}
}

func D1(key string, value interface{}) primitive.D {
	return primitive.D{
		E(key, value),
	}
}

func D2(key1 string, value1 interface{}, key2 string, value2 interface{}) primitive.D {
	return primitive.D{
		E(key1, value1),
		E(key2, value2),
	}
}

func D3(key1 string, value1 interface{}, key2 string, value2 interface{}, key3 string, value3 interface{}) primitive.D {
	return primitive.D{
		E(key1, value1),
		E(key2, value2),
		E(key3, value3),
	}
}
