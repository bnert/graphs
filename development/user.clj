(ns user)

(require '[graphs.core :as graph] :reload)



; A 'general purpose' graph can be represented
; via a hash map, assuming that each node in the graph is unique and uniquely
; key-able.
;
; The semantics of a node/edge relationship can be:
;
;  {"node-one" {"node-one" "edge-value"
;               "node-two" "edge-value"}
;   "node-two" {"node-three" "edge-value"}}
;
; which (in a textual representation, looks like:
;
;             |--- edge-value ---(node-two)--- edge-value ---(node-three)
;  (node-one)-|
;      |      |
;      |      |
;      \______/
;
(def unidirectional-non-cyclical
  {"a" {"b" nil
        "c" 5}
   "b" {}
   "c" {"d" 3}
   "d" {"e" 4}
   "e" {"f" 10}
   "f" {"e" 10}
   ; lil' orphan z
   "z" nil})

(def unidirectional
  {"a" #{"b" "c"}
   "b" {}
   "c" #{"a" "f"}
   "d" #{"e"}
   "e" {}
   "f" #{"e" "a" "c"}

   ; lil' orphan z
   "z" {}
   })

(def units
  (graph/bidi
    {"m"  {"ft" 3.28084}
     "ft" {"in" 12}
     "in" {}}
    #(/ 1 %)))

(comment
  (get-in unidirectional ["a" "c"])

  (let [pw (comp (partial graph/path->weights units)
                 (partial graph/path units))]
    (reduce * 10 (pw "m" "in")))

  (reduce * 10 (graph/weights units "m" "in"))

  ; outgoing node/edge values for "a"
  (graph/outgoing? unidirectional "a")
  (graph/outgoing? unidirectional "c")
  (graph/outgoing? unidirectional "e")

  (graph/incoming? unidirectional "a")

  (time
    (graph/cyclical? unidirectional))

  (time
    (graph/cyclical? unidirectional-non-cyclical))

  (graph/dfs unidirectional-non-cyclical "a" "f")
  (graph/bfs unidirectional-non-cyclical "a" "f")

  (graph/dfs units "m" "in")
  (graph/bfs units "m" "in")

  (graph/path unidirectional "a" "f")
  (graph/path unidirectional-non-cyclical "a" "f")
  (graph/path unidirectional "a" "z")
  (graph/path unidirectional "a" "b")
  (graph/path unidirectional "a" "a")

  (graph/path->weights
    unidirectional
    (graph/path unidirectional "a" "f"))

  (let [g unidirectional-non-cyclical]
    (graph/path->weights g
                         (graph/path g "a" "f")))

  (let [q (graph/q! ["a" "b" "c" "d"])]
    (println (peek q))
    (println (pop q))
    (println (apply conj (pop q) ["1" "2" "3"]))
    (println (into (pop q) ["1" "2" "3"]))
    (println q))



  ; window iter
  (let [path        ["a" "b" "c"]
        window-size 2]
    (loop [path' path
           steps []]
      (if (< (count path') window-size)
        steps
        (recur (vec (rest path'))
               (conj steps (subvec path' 0 window-size))))))




)
