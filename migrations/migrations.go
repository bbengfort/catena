package migrations

// Code generated by go generate; DO NOT EDIT.

func init() {
	migrations = make([]Migration, 0, 2)
	local(0, "migrations schema", "0000_migrations_schema.sql", []byte{67, 82, 69, 65, 84, 69, 32, 84, 65, 66, 76, 69, 32, 73, 70, 32, 78, 79, 84, 32, 69, 88, 73, 83, 84, 83, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 32, 40, 32, 34, 114, 101, 118, 105, 115, 105, 111, 110, 34, 32, 105, 110, 116, 101, 103, 101, 114, 32, 78, 79, 84, 32, 78, 85, 76, 76, 44, 32, 34, 110, 97, 109, 101, 34, 32, 118, 97, 114, 99, 104, 97, 114, 40, 49, 50, 56, 41, 32, 78, 79, 84, 32, 78, 85, 76, 76, 44, 32, 34, 97, 99, 116, 105, 118, 101, 34, 32, 98, 111, 111, 108, 101, 97, 110, 32, 78, 79, 84, 32, 78, 85, 76, 76, 32, 68, 69, 70, 65, 85, 76, 84, 32, 102, 97, 108, 115, 101, 44, 32, 34, 97, 112, 112, 108, 105, 101, 100, 34, 32, 84, 73, 77, 69, 83, 84, 65, 77, 80, 32, 87, 73, 84, 72, 32, 84, 73, 77, 69, 32, 90, 79, 78, 69, 44, 32, 34, 99, 114, 101, 97, 116, 101, 100, 34, 32, 84, 73, 77, 69, 83, 84, 65, 77, 80, 32, 87, 73, 84, 72, 32, 84, 73, 77, 69, 32, 90, 79, 78, 69, 32, 78, 79, 84, 32, 78, 85, 76, 76, 32, 68, 69, 70, 65, 85, 76, 84, 32, 67, 85, 82, 82, 69, 78, 84, 95, 84, 73, 77, 69, 83, 84, 65, 77, 80, 44, 32, 80, 82, 73, 77, 65, 82, 89, 32, 75, 69, 89, 32, 40, 34, 114, 101, 118, 105, 115, 105, 111, 110, 34, 41, 44, 32, 41, 32, 87, 73, 84, 72, 79, 85, 84, 32, 79, 73, 68, 83, 59, 32, 67, 79, 77, 77, 69, 78, 84, 32, 79, 78, 32, 84, 65, 66, 76, 69, 32, 34, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 34, 32, 73, 83, 32, 39, 77, 97, 110, 97, 103, 101, 115, 32, 116, 104, 101, 32, 115, 116, 97, 116, 101, 32, 111, 102, 32, 100, 97, 116, 97, 98, 97, 115, 101, 32, 98, 121, 32, 101, 110, 97, 98, 108, 105, 110, 103, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 32, 97, 110, 100, 32, 114, 111, 108, 108, 98, 97, 99, 107, 115, 39, 59, 32, 67, 79, 77, 77, 69, 78, 84, 32, 79, 78, 32, 67, 79, 76, 85, 77, 78, 32, 34, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 34, 46, 34, 114, 101, 118, 105, 115, 105, 111, 110, 34, 32, 73, 83, 32, 39, 84, 104, 101, 32, 114, 101, 118, 105, 115, 105, 111, 110, 32, 105, 100, 32, 112, 97, 114, 115, 101, 100, 32, 102, 114, 111, 109, 32, 116, 104, 101, 32, 102, 105, 108, 101, 110, 97, 109, 101, 32, 111, 102, 32, 116, 104, 101, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 39, 59, 32, 67, 79, 77, 77, 69, 78, 84, 32, 79, 78, 32, 67, 79, 76, 85, 77, 78, 32, 34, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 34, 46, 34, 110, 97, 109, 101, 34, 32, 73, 83, 32, 39, 84, 104, 101, 32, 110, 97, 109, 101, 32, 111, 102, 32, 116, 104, 101, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 32, 112, 97, 114, 115, 101, 100, 32, 102, 114, 111, 109, 32, 116, 104, 101, 32, 102, 105, 108, 101, 110, 97, 109, 101, 32, 111, 102, 32, 116, 104, 101, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 39, 59, 32, 67, 79, 77, 77, 69, 78, 84, 32, 79, 78, 32, 67, 79, 76, 85, 77, 78, 32, 34, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 34, 46, 34, 97, 99, 116, 105, 118, 101, 34, 32, 73, 83, 32, 39, 73, 102, 32, 116, 104, 101, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 32, 104, 97, 115, 32, 98, 101, 101, 110, 32, 97, 112, 112, 108, 105, 101, 100, 44, 32, 115, 101, 116, 32, 116, 111, 32, 102, 97, 108, 115, 101, 32, 111, 110, 32, 114, 111, 108, 108, 98, 97, 99, 107, 115, 32, 111, 114, 32, 105, 102, 32, 110, 111, 116, 32, 97, 112, 112, 108, 105, 101, 100, 39, 59, 32, 67, 79, 77, 77, 69, 78, 84, 32, 79, 78, 32, 67, 79, 76, 85, 77, 78, 32, 34, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 34, 46, 34, 97, 112, 112, 108, 105, 101, 100, 34, 32, 73, 83, 32, 39, 84, 105, 109, 101, 115, 116, 97, 109, 112, 32, 119, 104, 101, 110, 32, 116, 104, 101, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 32, 119, 97, 115, 32, 97, 112, 112, 108, 105, 101, 100, 44, 32, 110, 117, 108, 108, 32, 105, 102, 32, 114, 111, 108, 108, 101, 100, 98, 97, 99, 107, 32, 111, 114, 32, 110, 111, 116, 32, 97, 112, 112, 108, 105, 101, 100, 39, 59, 32, 67, 79, 77, 77, 69, 78, 84, 32, 79, 78, 32, 67, 79, 76, 85, 77, 78, 32, 34, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 34, 46, 34, 99, 114, 101, 97, 116, 101, 100, 34, 32, 73, 83, 32, 39, 84, 105, 109, 101, 115, 116, 97, 109, 112, 32, 119, 104, 101, 110, 32, 116, 104, 101, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 32, 119, 97, 115, 32, 99, 114, 101, 97, 116, 101, 100, 39, 59, 32}, []byte{68, 82, 79, 80, 32, 84, 65, 66, 76, 69, 32, 73, 70, 32, 69, 88, 73, 83, 84, 83, 32, 109, 105, 103, 114, 97, 116, 105, 111, 110, 115, 32, 67, 65, 83, 67, 65, 68, 69, 59, 32})
	local(1, "foo", "0001_foo.sql", []byte{}, []byte{})
}
