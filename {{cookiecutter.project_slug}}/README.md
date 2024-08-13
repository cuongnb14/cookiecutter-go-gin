# Migrate DB guide

## Setup
- `brew install ariga/tap/atlas`
- Config atlas in `atlas.hcl`
- Config gorm loader in `cmd/atlasloader/main.go`

## Workflow
- Create migration file: `atlas migrate diff --env gorm`
- Apply new migration file: `atlas migrate apply --url "postgres://dev:dev@devhost:5432/{{ cookiecutter.project_slug }}?sslmode=disable" --dir file://migrations`

## Create migrate seed data
- Create empty migration file: `atlas migrate new <migrate_name>`
- Update migration file
- Re-calculate hash: `atlas migrate hash`

# Utils
```sh
# Create model with repository and service
inv gen-model <ModelName>

# Format env: Update file env/template.env first
inv format-env <stage | all>
```
