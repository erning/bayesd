Build
```
$ ge go get https://github.com/erning/bayesd
$ cd bayesd
$ ge go build
$ # GOOS=linux GOARCH=amd64 ge go build
```

New
```
$ bayesd -new good bad
```

Learn
```
$ curl -X POST "http://127.0.0.1:8080/learn/good" -d '["tall", "rich", "handsome"]'
$ curl -X POST "http://127.0.0.1:8080/learn/bad" -d '["poor", "smelly", "ugly"]'
```

Guess
```
$ curl -X POST "http://127.0.0.1:8080/guess" -d '["tall", "girl"]' | json_pp

{
   "strict" : true,
   "scores" : {
      "bad" : 2.99999999991e-11,
      "good" : 0.99999999997
   },
   "likely" : "good"
}
```

Load from data file
```
$ bayesd
```
