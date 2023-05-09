package model

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type Filter struct {
	Key   string
	Value string
}
type Element struct {
	Name   string
	Filter Filter
}

type YamlPath []Element

func (p YamlPath) YamlPath() string {
	if len(p) == 0 {
		return ""
	}

	res := p[0].Name
	for i := 1; i < len(p); i++ {
		res += "." + p[i].Name
	}

	return res
}

func (p YamlPath) FilterFor(tag string) Filter {
	for _, elem := range p {
		if elem.Name == tag {
			return elem.Filter
		}
	}

	return Filter{}
}

func (f Filter) Search(nodes []*yaml.Node) bool {
	for i := 0; i < len(nodes); i += 2 {
		if nodes[i].Value == f.Key && nodes[i+1].Value == f.Value {
			return true
		}
	}

	return false
}

func ParseYamlPath(val string) YamlPath {
	var res YamlPath
	splitted := strings.Split(val, ".")
	for _, s := range splitted {
		if strings.Contains(s, "[") {
			elem, filterVal, _ := strings.Cut(s, "[")
			filterVal = strings.TrimSuffix(filterVal, "]")
			keyVal := strings.Split(filterVal, "=")
			res = append(res, Element{
				Name: elem,
				Filter: Filter{
					Key:   keyVal[0],
					Value: keyVal[1],
				},
			})
		} else {
			res = append(res, Element{Name: s})
		}
	}

	return res
}
