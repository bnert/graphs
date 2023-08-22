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
        "c" nil}
   "b" {}
   "c" {"d" nil}
   "d" {"e" nil}
   "e" {"f" nil}
   "f" {"e" nil}
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

(comment

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

  (graph/path unidirectional "a" "f")
  (graph/path unidirectional-non-cyclical "a" "f")
  (graph/path unidirectional "a" "z")
  (graph/path unidirectional "a" "b")
  (graph/path unidirectional "a" "a")

  (let [q (graph/q! ["a" "b" "c" "d"])]
    (println (peek q))
    (println (pop q))
    (println (apply conj (pop q) ["1" "2" "3"]))
    (println (into (pop q) ["1" "2" "3"]))
    (println q))
)