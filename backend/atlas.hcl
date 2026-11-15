data "external_schema" "gorm" {
  program = ["go", "run", "-mod=mod", "./cmd/atlasloader"]
}

env "local" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/16/dev?search_path=public"
  url = getenv("DATABASE_URL_DIRECT")
  revisions_schema = "public"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
