Install
```
$ ge go install https://github.com/erning/bayesd
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
$ curl -X POST "http://127.0.0.1:8080/guess" -d '["tall", "girl"]'

{
   "strict" : true,
   "likely" : "bad",
   "scores" : {
      "good" : 2.99999999991e-11,
      "bad" : 0.99999999997
   }
}
```

Load from data file
```
$ bayesd
```
