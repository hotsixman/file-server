# Multi-Protocol File Server (WebDAV & FTP)

이 프로젝트는 Go 언어로 구현된 고성능 멀티 프로토콜 파일 서버입니다. 하나의 저장소(
ootDir)를 WebDAV와 FTP 프로토콜을 통해 동시에 접근하고 관리할 수 있습니다.

## 🚀 주요 특징

- **멀티 프로토콜 지원**: WebDAV 및 FTP 지원.
- **통합 권한 관리**: auth.json을 통한 사용자별 비밀번호 및 경로 기반 접근 제어 (Match/Prefix).
- **공유 락 시스템**: WebDAV와 FTP 간의 동시 수정 충돌을 방지하기 위한 통합 LockSystem 적용.
- **고성능 엔진**: 
  - WebDAV: Gin Gonic + golang.org/x/net/webdav
  - FTP: ftpserverlib 기반 커스텀 드라이버 구현

## 📁 프로젝트 구조

```	ext
├── main.go             # 서버 진입점 (WebDAV/FTP 동시 실행)
├── auth.json           # 사용자 인증 및 권한 설정 파일
├── spec/               # WebDAV 상세 명세서
└── module/
    ├── auth/           # 인증 및 권한 검증 로직
    ├── webdav/         # WebDAV 서버 구현 (Gin 연동)
    ├── ftp/            # FTP 서버 구현 (Main/Client 드라이버)
    └── LS/             # 통합 LockSystem (MemoryLS 등)
```

## 🛠 설정 방법

### 1. 사용자 권한 설정 (auth.json)
사용자별로 접근 가능한 경로를 설정할 수 있습니다.

`json
{
  "user1": {
    "password": "pass123",
    "match": ["/"],
    "prefix": ["/data/"]
  }
}
`
- match: 정확히 일치하는 경로만 허용.
- prefix: 해당 경로로 시작하는 모든 하위 리소스 허용.

### 2. 실행
```sh
go mod tidy
go run main.go
```

## 📡 프로토콜 접속 정보

### WebDAV
- **URL**: http://localhost:3000/dav
- **Auth**: Basic Auth (ID/PW)
- **지원 메서드**: PROPFIND, MKCOL, LOCK, UNLOCK, GET, PUT 등 모든 필수 메서드.

### FTP
- **Address**: localhost:21
- **Auth**: User/Password
- **Mode**: Passive/Active 지원

## 📝 개발 가이드
상세한 WebDAV 동작 방식은 spec/webdav.md 문서를 참고하세요.
