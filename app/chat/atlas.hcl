data "external_schema" "gorm" {
  program = [
    "env",
    "LOG_LEVEL=ERROR",
    "go",
    "run",
    "./cmd/migrate",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
    exclude = ["public.*[type=function]", "public.*[type=table].*[type=trigger]"]
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}