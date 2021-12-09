package server

import "testing"

//
//func TestServer(t *testing.T) {
//	Server("8080")
//}
//
//func TestClient(t *testing.T) {
//	Client(8080)
//}
func BenchmarkClient(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Client(8080)
	}
}
