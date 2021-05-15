First Demo :
---
A basic unary operation between client and server.

Run `go run server/main.go`

---
Second Demo :
---
A git submodule way to store protobufs files so that you dont have to maintain multiple versions.

Show `protos` folder and `.gitmodules` file

Add another git repo as a submodule to current repo.

- `git submodule add git@github.com:mahendrabagul/protobufs.git protos`

Update existing git submodule.

- `git submodule update --remote`

---
Third Demo :
---
Running golang-grpc-server using server certificates.

Show `main.go` file and compare current changes with the earlier commit

---
Fourth Demo :
---
Enable golang-grpc-server with mTLS settings.

- Show `main.go` file and compare current changes with the earlier commit

---
Fifth Demo :
---
Add a check for verifying client certificates in golang-grpc-server with mTLS settings.

- Show `main.go` file and compare current changes with the earlier commit

---
