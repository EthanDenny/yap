# yap

A LISP written in Go for fun and learning.

```clojure
(def fizzbuzz [max n] {
    (= n (if-else (is-none n) 1 n))

    (= str (if-else (== (% n 3) 0) "Fizz" ""))
    (= str (if-else (== (% n 5) 0) (++ str "Buzz") str))
    (print (if-else (== str "") n str))

    (if (< n max) {
        (print ", ")
        (fizzbuzz max (+ n 1))
    })
})

(fizzbuzz 100)
```
