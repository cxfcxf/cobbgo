package main

import (
  "encoding/json"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "github.com/martini-contrib/binding"
  "io/ioutil"
)

type Profile struct {
  Hostname  string  `form:"hostname"`
  Macaddr  string `form:"macaddr"`
  Ksprofile  string  `form:"ksprofile"`
}

func addprofile(config map[string]map[string]string, addp Profile) bool {
    config[addp.Hostname] = make(map[string]string)
    config[addp.Hostname]["macaddr"] = addp.Macaddr
    config[addp.Hostname]["ksprofile"] = addp.Ksprofile
    c, _ := json.Marshal(config)
    err := ioutil.WriteFile("config.json", c, 0644)
    if err != nil {
        return false
    }
    return true
}

//func delprofile(config map[string]map[string]string, delp Delprofile) bool {
//    delete(config, delp.Hostname)
//    c, _ := json.Marshal(config)
//    err := ioutil.WriteFile("config.json", c, 0644)
//    if err != nil {
//        return false
//    }
//    return true
//}

func main() {
  conf, _ := ioutil.ReadFile("config.json")

  var config map[string]map[string]string

  json.Unmarshal(conf, &config)

  m := martini.Classic()
  m.Use(render.Renderer())

  m.Get("/cobbgo", func(r render.Render) {
    r.HTML(200, "index", config)
  })

//  m.Delete("/cobbgo/:name", func(params martini.Params, r render.Render) {
//    if delprofile(config, params["name"]) {
//        r.HTML(200, "index", config)
//    } else {
//        r.HTML(200, "index", config)
//    }
//  })

  m.Get("/cobbgo/add", func(r render.Render) {
    r.HTML(200, "add", "")
  })

  m.Post("/cobbgo/add", binding.Bind(Profile{}), func(addp Profile, r render.Render) {
    if addprofile(config, addp) {
        r.HTML(200, "index", config)
    } else {
        r.HTML(200, "index", config)
    }
  })

  m.Run()
}
