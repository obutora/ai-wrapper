// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/obutora/go_graphql_template/generated/ent/dot"
)

// CollectFields tells the query-builder to eagerly load connected nodes by resolver context.
func (d *DotQuery) CollectFields(ctx context.Context, satisfies ...string) (*DotQuery, error) {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return d, nil
	}
	if err := d.collectField(ctx, false, graphql.GetOperationContext(ctx), fc.Field, nil, satisfies...); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *DotQuery) collectField(ctx context.Context, oneNode bool, opCtx *graphql.OperationContext, collected graphql.CollectedField, path []string, satisfies ...string) error {
	path = append([]string(nil), path...)
	var (
		unknownSeen    bool
		fieldSeen      = make(map[string]struct{}, len(dot.Columns))
		selectedFields = []string{dot.FieldID}
	)
	for _, field := range graphql.CollectFields(opCtx, collected.Selections, satisfies) {
		switch field.Name {
		case "uid":
			if _, ok := fieldSeen[dot.FieldUID]; !ok {
				selectedFields = append(selectedFields, dot.FieldUID)
				fieldSeen[dot.FieldUID] = struct{}{}
			}
		case "title":
			if _, ok := fieldSeen[dot.FieldTitle]; !ok {
				selectedFields = append(selectedFields, dot.FieldTitle)
				fieldSeen[dot.FieldTitle] = struct{}{}
			}
		case "createdAt":
			if _, ok := fieldSeen[dot.FieldCreatedAt]; !ok {
				selectedFields = append(selectedFields, dot.FieldCreatedAt)
				fieldSeen[dot.FieldCreatedAt] = struct{}{}
			}
		case "updatedAt":
			if _, ok := fieldSeen[dot.FieldUpdatedAt]; !ok {
				selectedFields = append(selectedFields, dot.FieldUpdatedAt)
				fieldSeen[dot.FieldUpdatedAt] = struct{}{}
			}
		case "id":
		case "__typename":
		default:
			unknownSeen = true
		}
	}
	if !unknownSeen {
		d.Select(selectedFields...)
	}
	return nil
}

type dotPaginateArgs struct {
	first, last   *int
	after, before *Cursor
	opts          []DotPaginateOption
}

func newDotPaginateArgs(rv map[string]any) *dotPaginateArgs {
	args := &dotPaginateArgs{}
	if rv == nil {
		return args
	}
	if v := rv[firstField]; v != nil {
		args.first = v.(*int)
	}
	if v := rv[lastField]; v != nil {
		args.last = v.(*int)
	}
	if v := rv[afterField]; v != nil {
		args.after = v.(*Cursor)
	}
	if v := rv[beforeField]; v != nil {
		args.before = v.(*Cursor)
	}
	if v, ok := rv[orderByField]; ok {
		switch v := v.(type) {
		case []*DotOrder:
			args.opts = append(args.opts, WithDotOrder(v))
		case []any:
			var orders []*DotOrder
			for i := range v {
				mv, ok := v[i].(map[string]any)
				if !ok {
					continue
				}
				var (
					err1, err2 error
					order      = &DotOrder{Field: &DotOrderField{}, Direction: entgql.OrderDirectionAsc}
				)
				if d, ok := mv[directionField]; ok {
					err1 = order.Direction.UnmarshalGQL(d)
				}
				if f, ok := mv[fieldField]; ok {
					err2 = order.Field.UnmarshalGQL(f)
				}
				if err1 == nil && err2 == nil {
					orders = append(orders, order)
				}
			}
			args.opts = append(args.opts, WithDotOrder(orders))
		}
	}
	if v, ok := rv[whereField].(*DotWhereInput); ok {
		args.opts = append(args.opts, WithDotFilter(v.Filter))
	}
	return args
}

const (
	afterField     = "after"
	firstField     = "first"
	beforeField    = "before"
	lastField      = "last"
	orderByField   = "orderBy"
	directionField = "direction"
	fieldField     = "field"
	whereField     = "where"
)

func fieldArgs(ctx context.Context, whereInput any, path ...string) map[string]any {
	field := collectedField(ctx, path...)
	if field == nil || field.Arguments == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	args := field.ArgumentMap(oc.Variables)
	return unmarshalArgs(ctx, whereInput, args)
}

// unmarshalArgs allows extracting the field arguments from their raw representation.
func unmarshalArgs(ctx context.Context, whereInput any, args map[string]any) map[string]any {
	for _, k := range []string{firstField, lastField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		i, err := graphql.UnmarshalInt(v)
		if err == nil {
			args[k] = &i
		}
	}
	for _, k := range []string{beforeField, afterField} {
		v, ok := args[k]
		if !ok {
			continue
		}
		c := &Cursor{}
		if c.UnmarshalGQL(v) == nil {
			args[k] = c
		}
	}
	if v, ok := args[whereField]; ok && whereInput != nil {
		if err := graphql.UnmarshalInputFromContext(ctx, v, whereInput); err == nil {
			args[whereField] = whereInput
		}
	}

	return args
}

// mayAddCondition appends another type condition to the satisfies list
// if it does not exist in the list.
func mayAddCondition(satisfies []string, typeCond []string) []string {
Cond:
	for _, c := range typeCond {
		for _, s := range satisfies {
			if c == s {
				continue Cond
			}
		}
		satisfies = append(satisfies, c)
	}
	return satisfies
}
