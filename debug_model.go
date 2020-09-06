package pgext

import (
	"fmt"
	"reflect"

	"github.com/go-pg/pg/v10/orm"
	"gopkg.in/yaml.v2"
)

type debugTable struct {
	TableName string `yaml:"table_name"`

	PKs       []debugColumn   `yaml:"pks"`
	Columns   []debugColumn   `yaml:"columns"`
	Relations []debugRelation `yaml:"relations"`
}

type debugColumn struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type debugRelation struct {
	Name      string   `yaml:"name"`
	JoinType  string   `yaml:"join_type"`
	JoinTable string   `yaml:"join_table"`
	JoinOn    []string `yaml:"join_on"`
}

func DebugModel(model interface{}) string {
	typ := reflect.TypeOf(model)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	table := orm.GetTable(typ)

	dt := new(debugTable)
	dt.TableName = table.Name

	for _, pk := range table.PKs {
		dt.PKs = append(dt.PKs, debugColumn{
			Name: pk.SQLName,
			Type: pk.SQLType,
		})
	}

	for _, f := range table.DataFields {
		dt.Columns = append(dt.Columns, debugColumn{
			Name: f.SQLName,
			Type: f.SQLType,
		})
	}

	for _, rel := range table.Relations {
		dt.Relations = append(dt.Relations, newDebugRelation(table, rel))
	}

	b, err := yaml.Marshal(dt)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func newDebugRelation(base *orm.Table, rel *orm.Relation) debugRelation {
	join := rel.JoinTable

	dr := debugRelation{}
	dr.Name = rel.Field.GoName
	dr.JoinTable = join.Name

	switch rel.Type {
	case orm.HasOneRelation:
		dr.JoinType = "has one"
		for i, fk := range rel.FKs {
			dr.JoinOn = append(dr.JoinOn, fmt.Sprintf(
				"%s.%s = %s.%s",
				join.Name, join.PKs[i].SQLName,
				base.Name, fk.SQLName,
			))
		}
	case orm.BelongsToRelation:
		dr.JoinType = "belongs to"
	case orm.HasManyRelation:
		dr.JoinType = "has many"
		for i, fk := range rel.FKs {
			dr.JoinOn = append(dr.JoinOn, fmt.Sprintf(
				"%s.%s = %s.%s",
				join.Name, fk.SQLName,
				base.Name, rel.FKValues[i].SQLName,
			))
		}
	case orm.Many2ManyRelation:
		dr.JoinType = "many to many"
		for i, fk := range rel.M2MBaseFKs {
			dr.JoinOn = append(dr.JoinOn, fmt.Sprintf(
				"%s.%s = %s.%s",
				rel.M2MTableAlias, fk,
				base.Name, base.PKs[i].SQLName,
			))
		}
		for i, fk := range rel.M2MJoinFKs {
			dr.JoinOn = append(dr.JoinOn, fmt.Sprintf(
				"%s.%s = %s.%s",
				rel.M2MTableAlias, fk,
				join.Name, join.PKs[i].SQLName,
			))
		}
	default:
		panic("not reached")
	}

	return dr
}
