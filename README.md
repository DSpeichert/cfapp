# cfapp

Build, Test, Use:
```bash
go get
go build
swagger generate spec -o ./swagger.json # Optionally: to rebuild swagger.json
./cfapp help
./cfapp -m 'user:password@tcp(localhost:3306)/cfapp?charset=utf8' -n http://requestb.in/idhere
# Open http://petstore.swagger.io/?url=http://localhost/swagger.json
```
