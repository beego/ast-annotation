package annotation

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestAnnotationVisitor_Visit(t *testing.T) {
	testCases := []struct {
		name     string
		src      string
		fileAns  []Annotation
		funcAns  []Annotation
		fieldAns []Annotation
		typeAns  []Annotation
	}{
		{
			name: "file case",
			src: `
// annotation go through the source code and extra the annotation
// @author Deng Ming
/* @multiple first line
second line*/
// @date 2022/04/02
package annotation
`,
			fileAns: []Annotation{
				{
					Key:   "author",
					Value: "Deng Ming",
				},
				{
					Key:   "multiple",
					Value: "first line\nsecond line",
				},
				{
					Key:   "date",
					Value: "2022/04/02",
				},
			},
		},
		{
			name: "type func",
			src: `
package annotation

type (
	// FuncType is a type
	// @author Deng Ming
	/* @multiple first line
	   second line*/
	// @date 2022/04/02
	FuncType func()
)
`,
			typeAns: []Annotation{
				{
					Key:   "author",
					Value: "Deng Ming",
				},
				{
					Key:   "multiple",
					Value: "first line\n\t   second line",
				},
				{
					Key:   "date",
					Value: "2022/04/02",
				},
			},
		},
		{
			name: "field",
			src: `
package annotation

type (
	// StructType is a test struct
	//
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/02
	StructType struct {
		// Public is a field
		// @type string
		Public string
	}
)
`,
			fieldAns: []Annotation{
				{
					Key:   "type",
					Value: "string",
				},
			},
			typeAns: []Annotation{
				{
					Key:   "author",
					Value: "Deng Ming",
				},
				{
					Key:   "multiple",
					Value: "first line\n\t   second line\n\t",
				},
				{
					Key:   "date",
					Value: "2022/04/02",
				},
			},
		},
		{
			name: "func",
			src: `
`,
			typeAns: []Annotation{
				{
					Key:   "author",
					Value: "Deng Ming",
				},
				{
					Key:   "multiple",
					Value: "first line\n\t   second line\n\t",
				},
				{
					Key:   "date",
					Value: "2022/04/02",
				},
			},
			funcAns: []Annotation{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet() // positions are relative to fset
			f, err := parser.ParseFile(fset, "src.go", tc.src, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			fv := &fileAnnotationVisitor{}
			fnv := &funcAnnotationVisitor{}
			fdv := &fieldAnnotationVisitor{}
			tv := &typeAnnotationVisitor{}
			visitor := &AnnotationVisitor{
				File:     fv,
				FuncDecl: fnv,
				Field:    fdv,
				Type:     tv,
			}
			ast.Walk(visitor, f)
			assertAnnotations(t, tc.funcAns, fnv.ans)
			assertAnnotations(t, tc.fileAns, fv.ans)
			assertAnnotations(t, tc.fieldAns, fdv.ans)
			assertAnnotations(t, tc.typeAns, tv.ans)
		})
	}
}

func assertAnnotations(t *testing.T, want []Annotation, dst Annotations) {
	if len(want) != len(dst.Ans) {
		t.Fatal()
	}
	for _, an := range want {
		val, ok := dst.Ans[an.Key]
		assert.True(t, ok)
		assert.Equal(t, an.Value, val.Value)
	}
}

type fieldAnnotationVisitor struct {
	ans Annotations
}

func (v *fieldAnnotationVisitor) Visit(field *ast.Field, ans Annotations) {
	v.ans = ans
}

type typeAnnotationVisitor struct {
	ans Annotations
}

func (v *typeAnnotationVisitor) Visit(s *ast.TypeSpec, ans Annotations) {
	v.ans = ans
}

type funcAnnotationVisitor struct {
	ans Annotations
}

func (v *funcAnnotationVisitor) Visit(fn *ast.FuncDecl, ans Annotations) {
	v.ans = ans
}

type fileAnnotationVisitor struct {
	ans Annotations
}

func (v *fileAnnotationVisitor) Visit(f *ast.File, ans Annotations) {
	v.ans = ans
}
