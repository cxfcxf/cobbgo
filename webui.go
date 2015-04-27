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

type Editprofile struct {
  Hostname string
  Info  map[string]string
}

func addprofile(config map[string]map[string]string, addp Profile) bool {
    if len(config) == 0 {
      config = make(map[string]map[string]string)
    }
    if len(config[addp.Hostname]) > 0 { return false }
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

func editprofile(config map[string]map[string]string, editp Profile) bool {
    config[editp.Hostname]["macaddr"] = editp.Macaddr
    config[editp.Hostname]["ksprofile"] = editp.Ksprofile
    c, _ := json.Marshal(config)
    err := ioutil.WriteFile("config.json", c, 0644)
    if err != nil {
      return false
    }
    return true
}


func delprofile(config map[string]map[string]string, delp string) bool {
    delete(config, delp)
    c, _ := json.Marshal(config)
    err := ioutil.WriteFile("config.json", c, 0644)
    if err != nil {
        return false
    }
    return true
}


func readconf() map[string]map[string]string{
  conf, _ := ioutil.ReadFile("config.json")

  var config map[string]map[string]string

  json.Unmarshal(conf, &config)

  return config
}


func main() {

  m := martini.Classic()
  m.Use(render.Renderer())

  m.Get("/cobbgo", func(r render.Render) {
    config := readconf()
    r.HTML(200, "index", config)
  })

  m.Post("/cobbgo/add", binding.Bind(Profile{}), func(addp Profile, r render.Render) {
    config := readconf()
    if addprofile(config, addp) {
        r.Redirect("/cobbgo", 302)
    } else {
        r.Redirect("/cobbgo", 302)
    }
  })

  m.Post("/cobbgo/delete/:hostname", func(params martini.Params, r render.Render){
   config := readconf()
   if delprofile(config, params["hostname"]) {
        r.Redirect("/cobbgo", 302)
    } else {
        r.Redirect("/cobbgo", 302)
    }
  })

  m.Get("/cobbgo/edit/:hostname", func(params martini.Params, r render.Render) {
    config := readconf()
    sw := true
    for k, _ := range config {
      if k == params["hostname"] {
        var profile Editprofile
        profile.Hostname = params["hostname"]
        profile.Info = config[params["hostname"]]
        r.HTML(200, "edit", profile)
        sw = false
      }
    }
    if sw {
      r.Redirect("/cobbgo", 302)
    }
  })

  m.Post("/cobbgo/edit/:hostname", binding.Bind(Profile{}), func(params martini.Params, editp Profile, r render.Render) {
    config := readconf()
    editp.Hostname = params["hostname"]
    if editprofile(config, editp) {
      r.Redirect("/cobbgo", 302)
    } else {
      r.Redirect("/cobbgo", 302)
    }
  })

  m.Get("/cobbgo/add", func(r render.Render) {
    r.HTML(200, "add", "")
  })

  m.Run()
}
