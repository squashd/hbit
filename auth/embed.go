package auth

import "embed"

//go:embed schemas/*.sql
var Migrations embed.FS
