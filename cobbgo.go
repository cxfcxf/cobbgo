package main

import (
    "strings"
    "regexp"
    "encoding/json"
    "io/ioutil"
    "text/template"
    "os"
)

type Profile struct {
    Macaddr     string  `json:"macaddr"`
    Ksprofile   string  `json:"ksprofile"`
}

type Mac struct {
    Kernel, Initrd, Ksfile string
}

type Kickstart struct {
    Version, Ondisk, Offdisk, Ipaddr, Nm, Gw, Hostname
}

func filesgen(loc string) {
    r, _ := regexp.Compile(":")
    config, err := ioutil.ReadFile("config.json")
    if err != nil {panic(err)}

    var servers map[string]*Profile
    json.Unmarshal(config, &servers)

    t := template.Must(template.ParseFiles("templates/mac.tmpl"))

    for _, s := range servers {
        //create files
        var m Mac
        m.Kernel = "http://cxfcxf.netdnasa16.netdna-cdn.com/vmlinuz"
        m.Initrd = "http://cxfcxf.netdnasa16.netdna-cdn.com/initrd.img"
        m.Ksfile = s.Ksprofile

        macname := strings.ToLower(s.Macaddr)
        macname = r.ReplaceAllString(macname, "-")

        fd, err := os.Create(loc+"/"+macname)
        if err != nil {panic(err)}
        defer fd.Close()

        err = t.Execute(fd, m)
        if err != nil {panic(err)}
    }
}

func main() {
    filesgen("/Users/xuefengchen/myrepos/cobbgo/tftpboot")
}
