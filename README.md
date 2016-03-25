curl -XPOST localhost:8088/api/v1/aaa
curl -XDELETE localhost:8088/api/v1/aaa
curl -XPOST localhost:8088/api/v1/aaa/bbb
curl -XDELETE localhost:8088/api/v1/aaa/bbb
curl -XPUT localhost:8088/api/v1/aaa/bbb -F  "file=@t.jpg”
curl -XPUT localhost:8088/api/v1/aaa/bbb -F  "file=@t.jpg" -F "file=@f.jpg”

 curl -XPOST localhost:8088/api/v1/aaa/bbb/tt.jpg -F  "file=@t.jpg"
curl -XDELETE localhost:8088/api/v1/aaa/bbb/tt.jpg 