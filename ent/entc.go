//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {

	ex, err := entgql.NewExtension(
		// Tell Ent to generate a GraphQL schema for
		// the Ent schema in a file named ent.graphql.
		entgql.WithWhereInputs(true),
		entgql.WithSchemaGenerator(),
		entgql.WithSchemaPath("./schema/ent.graphql"),
		entgql.WithConfigPath("./gqlgen.yml"),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	if err := entc.Generate("./schema", &gen.Config{
		Target:  "./generated/ent",
		Package: "github.com/obutora/go_graphql_template/generated/ent",
		Features: []gen.Feature{
			gen.FeatureExecQuery,
			gen.FeatureUpsert,
			gen.FeatureModifier,
		},
	}, entc.Extensions(ex), entc.FeatureNames("intercept", "schema/snapshot")); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
