package evaluator

import (
	"github.com/davidsbond/language/ast"
	"github.com/davidsbond/language/object"
)

func blockStatement(block *ast.BlockStatement, scope *object.Scope) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Evaluate(statement, scope)

		if result != nil {
			rt := result.Type()
			if rt == object.TypeReturnValue || rt == object.TypeError {
				return result
			}
		}
	}

	return result
}
