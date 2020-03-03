//config.go
package config

import(
    "flag"
)

type Config struct {
    IP      string
    Port    int
    DB_path string
    Template_path string
}

func ParseFlags() Config {
    // params
    var ip = flag.String("ip", "0.0.0.0", "IP address to listen on in the form 1.2.3.4")
    var port = flag.Int("port", 80, "port to listen on")
    var db_path = flag.String("db", "db", "Absolute path to the jimbob database")
    var template_path = flag.String("t", "templates", "Absolute path to the templates")
    flag.Parse()

    return Config{*ip, *port, *db_path, *template_path}
}
