package util

import (
	"fmt"
	"reflect"

	"go.yaml.in/yaml/v3"
)

const (
	yamlTagTimestamp = "!!timestamp"
	yamlTagMerge     = "!!merge"
)

// DecodeYAML decodes YAML source (YAML 1.2, thus yes, no, on and off are
// strings) into generic values:
// - mappings are decoded as map[interface{}]interface{},
// - sequences are decoded as []interface{},
// - timestamps are kept as strings,
// - other scalars are decoded as string, int, float64, bool or nil.
// Anchors, aliases and merge keys (<<) are supported.
// - source: the YAML source to decode
// Return: decoded value and an error if something went wrong
func DecodeYAML(source []byte) (interface{}, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(source, &node); err != nil {
		return nil, err
	}
	return nodeToValue(&node)
}

// DecodeYAMLObject decodes YAML source that must be a mapping with string keys
// (or an empty document):
// - source: the YAML source to decode
// Return: decoded object and an error if something went wrong
func DecodeYAMLObject(source []byte) (Object, error) {
	value, err := DecodeYAML(source)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return Object{}, nil
	}
	return NewObject(value)
}

// UnmarshalYAML decodes YAML source into a structure:
// - source: the YAML source to decode
// - out: pointer to the structure to fill
// Return: an error if something went wrong
func UnmarshalYAML(source []byte, out interface{}) error {
	return yaml.Unmarshal(source, out)
}

// EncodeYAML encodes given value in YAML:
// - value: the value to encode
// Return: YAML encoded value and an error if something went wrong
func EncodeYAML(value interface{}) ([]byte, error) {
	return yaml.Marshal(value)
}

func nodeToValue(node *yaml.Node) (interface{}, error) {
	switch node.Kind {
	case 0:
		// empty source
		return nil, nil
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return nil, nil
		}
		return nodeToValue(node.Content[0])
	case yaml.AliasNode:
		return nodeToValue(node.Alias)
	case yaml.ScalarNode:
		if node.Tag == yamlTagTimestamp {
			return node.Value, nil
		}
		var value interface{}
		if err := node.Decode(&value); err != nil {
			return nil, err
		}
		return value, nil
	case yaml.SequenceNode:
		list := make([]interface{}, 0, len(node.Content))
		for _, child := range node.Content {
			value, err := nodeToValue(child)
			if err != nil {
				return nil, err
			}
			list = append(list, value)
		}
		return list, nil
	case yaml.MappingNode:
		return mappingToValue(node)
	}
	return nil, fmt.Errorf("unsupported YAML node at line %d", node.Line)
}

func mappingToValue(node *yaml.Node) (map[interface{}]interface{}, error) {
	result := make(map[interface{}]interface{})
	// merged mappings first, so that explicit keys override them
	for i := 0; i+1 < len(node.Content); i += 2 {
		key, value := node.Content[i], node.Content[i+1]
		if key.Tag != yamlTagMerge {
			continue
		}
		if err := mergeMapping(result, value); err != nil {
			return nil, err
		}
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		keyNode, valueNode := node.Content[i], node.Content[i+1]
		if keyNode.Tag == yamlTagMerge {
			continue
		}
		key, err := nodeToValue(keyNode)
		if err != nil {
			return nil, err
		}
		if key != nil && !reflect.TypeOf(key).Comparable() {
			return nil, fmt.Errorf("invalid map key at line %d", keyNode.Line)
		}
		value, err := nodeToValue(valueNode)
		if err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}

// mergeMapping merges mapping(s) of a merge key (<<) into result. Earlier
// mappings of a list have precedence over later ones.
func mergeMapping(result map[interface{}]interface{}, node *yaml.Node) error {
	if node.Kind == yaml.AliasNode {
		node = node.Alias
	}
	var nodes []*yaml.Node
	switch node.Kind {
	case yaml.MappingNode:
		nodes = []*yaml.Node{node}
	case yaml.SequenceNode:
		for i := len(node.Content) - 1; i >= 0; i-- {
			nodes = append(nodes, node.Content[i])
		}
	default:
		return fmt.Errorf("merge key value must be a mapping at line %d", node.Line)
	}
	for _, n := range nodes {
		if n.Kind == yaml.AliasNode {
			n = n.Alias
		}
		if n.Kind != yaml.MappingNode {
			return fmt.Errorf("merge key value must be a mapping at line %d", n.Line)
		}
		mapping, err := mappingToValue(n)
		if err != nil {
			return err
		}
		for key, value := range mapping {
			result[key] = value
		}
	}
	return nil
}
