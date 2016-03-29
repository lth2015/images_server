curl -XGET localhost:8088/version
curl -XGET localhost:8088/healthz
curl -XGET localhost:8088/api/v1/accounts

curl -XPOST localhost:8088/api/v1/accounts/aaa
curl -XDELETE localhost:8088/api/v1/accounts/aaa
curl -XPOST localhost:8088/api/v1/accounts/aaa/containers/bbb
curl -XDELETE localhost:8088/api/v1/accounts/aaa/containers/bbb

curl -XPUT localhost:8088/api/v1/accounts/aaa/containers/bbb -F  "file=@t.png”
curl -XPUT localhost:8088/api/v1/accounts/aaa/containers/bbb -F  "file=@t.png" -F "file=@f.png”
curl -XDELETE localhost:8088/api/v1/accounts/aaa/containers/bbb

curl -XPOST localhost:8088/api/v1/accounts/aaa/containers/bbb/buckets/tt.png -F  "file=@t.png"
curl -XDELETE localhost:8088/api/v1/accounts/aaa/containers/bbb/buckets/tt.png 
curl -XPOST localhost:8088/api/v1/accounts/aaa/containers/bbb/buckets/tt.png -F  "file=@t.png"
curl -XGET localhost:8088/api/v1/accounts/aaa/containers/bbb/buckets/tt.png 
