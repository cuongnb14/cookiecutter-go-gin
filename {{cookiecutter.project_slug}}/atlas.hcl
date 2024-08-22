data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas_loader",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  # dev = "docker://postgres/14.1/dev"
  dev = "postgres://dev:dev@devhost:5432/atlas?sslmode=disable"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}