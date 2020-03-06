## Thử làm một cái URL shorten đơn giản bằng Golang

### Demo
Truy cập trang: [http://short-moe.xyz/](http://short-moe.xyz/ "http://short-moe.xyz/")

### Các bước chạy
- Yêu cầu máy tính cài Go version 1.11 trở lên
- Clone repo về máy
- Chạy một database Postgres. Tốt nhất là dùng Docker
- Tạo file **config.default.json** có dạng như sau:
```json
{
    "database": {
        "user": "tên user truy cập database Postgres",
        "password": "password truy cập postgres",
        "database": "tên database của postgres",
        "address": "địa chỉ chạy database postgres"
    }
}
```
- **go run main.go**

### Các API

1. GET /
- Dùng để render giao diện trang chủ: Trả về file index.html trong thư mục view

2. POST /shorten
- Đọc file JSON gửi lên. Cấu trúc của file JSON như sau:
```go
type ShortenReq struct {
	Path string
}
```
- Tạo một biến có kiểu dữ liệu như sau:
```go
type Url struct {
	Id         int         // ID tự tăng lưu trong database
    OriginUrl  string      // URL dài ngoằng do client gửi lên
	ShortenUrl string      // URL ngắn được tạo ra
}
```
- Biến Url sẽ được INSERT vào Postgres database

3. GET /{id}
- Lấy id client gửi lên
- Tìm trong Postgres database bản ghi nào có ShortenUrl == id, từ đó truy ra được OriginUrl tương ứng
- Redirect về OriginUrl

### Tất nhiên là sẽ có nhiều cách khác nhau. Lúc nào rảnh mình sẽ thử làm cách khác xem thế nào :D