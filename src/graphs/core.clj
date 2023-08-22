(ns graphs.core)

(defprotocol Outgoing
  (outgoing? [this node]
    "Outgoing node/edges for specified node.
    Returns nil if node does not exist in the graph"))

(defprotocol Incoming
  (incoming? [this node]
    "Incoming node/enges for specified node.
    Returns nil if node doesn ot exist in the graph"))

(defprotocol Cyclical
  (cyclical?
    [this]
    #_[this root-node]
    "Set of cycles in the graph. Uses depth first approach."))

(defprotocol DepthFirstSearch
  (dfs [this start-node pred]
    "Searches a graph, via depth first search, halting on a matched predicate"))

(defprotocol BreadthFirstSearch
  (bfs [this start-node pred]
    "Searches a graph, via breadth first search, halting on a matched predicate."))

(defprotocol Paths
  (path [this start-node end-node]
    "Returns path of nodes. Empty seq if no path.")
  (path->weights [this path]
    "Returns weights for a path.")
  (weights [this start-node end-node]
    "List of weights for a path, if path exists"))

(defn unweighted-seq->unweighted-map [nodes]
  (cond
    (some #(%1 nodes) [vector? set? seq?])
      (reduce #(assoc %1 %2 nil) {} nodes)
    (string? nodes)
      {nodes nil}
    :else
      (or nodes {})))

(defn map-graph-cyclical? [m node]
  ; uses a dfs approach
  (when (get m node)
    (loop [queue   (keys (outgoing? m node))
           visited #{node}]
      (cond
        (contains? visited (first queue))
          true
        (not (seq queue))
          false
        :else
          (let [[head & tail] queue]
            (recur (concat tail (keys (outgoing? m head)))
                   (conj visited head)))))))

(defmethod print-method clojure.lang.PersistentQueue [q, w] ; Overload the printer for queues so they look like fish
  (print-method '<- w)
  (print-method (seq q) w)
  (print-method '-< w))

(defn q!
  ([]
   (q! []))
  ([coll]
   (into clojure.lang.PersistentQueue/EMPTY coll)))

(defn bfs* [m start-node match?]
  (loop [queue   (q! [start-node])
         paths   {start-node nil}
         visited []]
    (cond
      (match? (last visited))
        (with-meta visited {:paths paths})
      (empty? queue)
        (with-meta [] {:paths {}})
      :else
        (let [h     (peek queue)
              vs    (set visited)
              nodes (keys (outgoing? m h))
              nv    (filter (complement
                              #(contains? vs %))
                            nodes)
              nq    (into (pop queue) nv)
              np    (reduce #(assoc %1 %2 h)
                            paths
                            nodes)]
          (if (contains? vs h)
            (recur nq np visited)
            (recur nq np (conj visited h)))))))

(defn dfs* [m start-node match?]
  (loop [stack   [start-node]
         paths   {start-node nil}
         visited []]
    (cond
      (match? (last visited))
        (with-meta visited {:paths paths})
      (empty? stack)
        (with-meta [] {:paths {}})
      :else
        (let [top   (peek stack)
              vs    (set visited)
              nodes (vec (keys (outgoing? m top)))
              nv    (vec (filter (complement
                                   #(contains? vs %))
                                 nodes))
              np    (reduce #(assoc %1 %2 top)
                            paths
                            nodes)
              nstack (into (pop stack) nv)]
          (if (contains? vs top)
            (recur nstack np top)
            (recur nstack np (conj visited top)))))))

(defn path->steps [path]
  (loop [path' path
         steps []]
    (if (< (count path') 2)
      steps
      (recur (vec (rest path'))
             (conj steps
                   (subvec path' 0 2))))))

(extend-type clojure.lang.PersistentArrayMap
  BreadthFirstSearch
  (bfs [m start-node end-node]
    (when (every? (partial contains? m) [start-node end-node])
      (bfs* m start-node (partial = end-node))))

  Cyclical
  (cyclical? [m]
    ; need to iterate over each node, given there could be a cycle
    ; which isn't connected to a node in a potentially disjoint graph
    (or (some #(map-graph-cyclical? m %1) (keys m)) false))

  DepthFirstSearch
  (dfs [m start-node end-node]
    (when (every? (partial contains? m) [start-node end-node])
      (dfs* m start-node (partial = end-node))))

  Incoming
  (incoming? [m node]
    (when (get m node)
      (loop [node-edge-set  (vec m)
             result         {}]
        (if-not (seq node-edge-set)
          result
          (let [[[node' edge-nodes] & r] node-edge-set
                edge-nodes (unweighted-seq->unweighted-map edge-nodes)]
            (if (contains? edge-nodes node)
              (recur r
                     (assoc result node' (get edge-nodes node)))
              (recur r
                     result)))))))

  Outgoing
  (outgoing? [m node]
    (when (contains? m node)
      (unweighted-seq->unweighted-map (get m node))))

  Paths
  (path [m start-node end-node]
    (let [r (bfs m start-node end-node)
          r' ((meta r) :paths)]
      (-> (bfs r' (last r) (first r))
          reverse
          vec)))
  (path->weights [m path]
    (println (path->steps path))
    (mapv
      (fn [[node-start node-end]]
        (let [v? (unweighted-seq->unweighted-map
                   (get m node-start))]
          (or (get v? node-end) 0)))
      (path->steps path)))
  (weights [m start-node end-node]
    (path->weights m (path m start-node end-node))))

(defn invert [m from-node node-edge-set edge-fn]
  (reduce-kv
    (fn [m' node edge]
      (update m' node (fnil assoc {}) from-node (edge-fn edge)))
    m
    node-edge-set))
(defmulti bidi
  (fn [t _] (type t)))

(defmethod bidi clojure.lang.PersistentArrayMap
 [m edge-fn]
 (reduce-kv
   (fn [m' k v]
     (invert m' k v edge-fn))
   m
   m))

