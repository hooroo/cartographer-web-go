package main

import (
  "io"
  "log"
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/eknkc/amber"
  //_ "github.com/lib/pq" //Ideally this will create a dependency on lib/pq
)

// Need some JSON types for unpacking data

type json_system_node map[string]string

type json_report map[string]string

type json_node struct {
  System json_system_node `json:system`
  Reports map[string]json_report `json:reports`
}


// Non-JSON types

type node_report struct {
  Identifier string
  Detail string
  Status string
}

type node struct {
  Name string
  Reports []node_report
  CssStatus string
}

func ShowNodes(w http.ResponseWriter, req *http.Request) {
  ns, _ := GetNodes()

  var node_data map[string][]node = make(map[string][]node, 0)
  var nodes []node = make([]node, 0, 1)

  for _, n := range ns {
    reports := make([]node_report, 0, 2)
    for name, report := range n.Reports {
      var kv_pairs string = ""
      for key, value := range report {
        kv_pairs += key + "=" + value + ", "
      }

      rep := node_report{name, kv_pairs, "success"}
      reports = append(reports, rep)
    }
    nodes = append(nodes, node{n.System["node_name"], reports, "success"})
  }

  node_data["Nodes"] = nodes

  compiler := amber.New()
  err := compiler.ParseFile("templates/index.amber")
  if err != nil {
    io.WriteString(w, "Amber parse error!\n" + err.Error())
    return
  }

  tpl, err := compiler.Compile()
  if err != nil {
    io.WriteString(w, "Amber compile error!\n" + err.Error())
    return
  }

  tpl.Execute(w, node_data)
} // End ShowNodes. ShowNodes should go somewhere else!

func GetNodes() (node_array []json_node, err error) {
  json_text := `[
    {
      "system": { "node_name": "host.local", "ip_address": "192.168.0.1" },
      "reports": {
        "unicorn": { "workers": "4", "memory_used": "4096" },
        "nginx": { "workers": "8", "memory_used": "8192" }
      }
    }
  ]`

  err = json.Unmarshal([]byte(json_text), &node_array)
  return
} //Interrim Fetcher for Node data. should rely on DB at some stage...


func main() {
  // Define router
  router := mux.NewRouter()

  // Define routes here
  router.Methods("GET").Subrouter().HandleFunc("/", ShowNodes)

  // Mount file servers for js and css
  router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

  // Start serving requests!
  err := http.ListenAndServe(":12345", router)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
