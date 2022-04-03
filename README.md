# ast-annotation

In some scenarios, we may want the "annotation" feature, specially for code generation scenario.

So we come up with the general annotation solution, including:
- how to define the annotation in your source code;
- some enhancement API for Go AST SDK which allows users to read the annotation from source code.

## Annotation Specification

An annotation must follow these two principles:
- The annotation must be part of the **document** of the specific AST node
- The annotation must follow one of these two patterns:
  - // @annotation_name annotation_value
  - /* @annotation_name annotation_value */

For example:
```go
package main
type (
  // MyType is test type
  // @description my test type
  /* @multi_line first line
     second line
  */
  MyType struct {

  }
)
```
There are two annotations:
- description => my test type
- multi_line => first line \n second line

> if you want to add annotation for type definitions, you cannot use the syntax "type A xxx"

## Q & A

### Why not fork the AST code and make the enhancement?
It's not worthy doing that.