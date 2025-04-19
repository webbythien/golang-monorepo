# Monorepo

Đây là một monorepo được quản lý bằng [Nx](https://nx.dev), chứa các ứng dụng và thư viện khác nhau.

## Cấu trúc dự án

Monorepo được tổ chức thành hai thư mục chính:

- `app/`: Chứa các ứng dụng (bao gồm iam và chat)
- `packages/`: Chứa các thư viện chia sẻ

## Bắt đầu

### Yêu cầu

- Node.js
- npm hoặc yarn
- [Nx CLI](https://nx.dev/getting-started/installation)

### Cài đặt

```sh
# Cài đặt các phụ thuộc
npm install
```

## Sử dụng Nx

### Chạy các tác vụ

Để chạy bất kỳ tác vụ nào với Nx, sử dụng lệnh:

```sh
npx nx <tác-vụ> <tên-dự-án>
```

Ví dụ:
```sh
npx nx build app.iam
npx nx test app.chat
```

### Chạy API

Để khởi động API, sử dụng lệnh:

```sh
npx nx run app.iam:serve
```

### Chạy tác vụ cho nhiều dự án

```sh
npx nx run-many -t <tác-vụ> -p <dự-án-1> <dự-án-2>
```

### Các tác vụ chính

- `serve`: Khởi động dịch vụ API
- `setup`: Thiết lập dự án
- `generate`: Tạo mã tự động
- `tidy`: Dọn dẹp và tổ chức mã
- `lint`: Kiểm tra lỗi mã
- `test`: Chạy kiểm thử đơn vị
- `e2e`: Chạy kiểm thử end-to-end
- `build`: Xây dựng dự án
- `dockerize`: Tạo container Docker

## Tạo một thư viện mới

```sh
npx nx g @nx/js:lib packages/tên-thư-viện --publishable --importPath=@monorepo/tên-thư-viện
```

## Phát hành

Để phiên bản và phát hành các dự án:

```sh
npx nx release
```

Thêm `--dry-run` để xem trước các thay đổi mà không thực sự phát hành.

## CI/CD

Dự án sử dụng các quy ước commit tuân theo chuẩn để quản lý phiên bản và tạo changelog tự động. Các loại commit được hỗ trợ:

- `docs`: Thay đổi tài liệu
- `chore`: Công việc bảo trì
- `refactor`: Tái cấu trúc mã
- `revamp`: Cải tiến module
- `build`: Thay đổi hệ thống build
- `ci`: Thay đổi CI

## Tài liệu tham khảo

- [Tài liệu Nx](https://nx.dev)
- [Quản lý phát hành với Nx](https://nx.dev/features/manage-releases)
- [Cộng đồng Nx](https://go.nx.dev/community)

## Giấy phép

MIT
