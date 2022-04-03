// Copyright 2022 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package annotation

import (
	"go/ast"
	"strings"
)

type Annotations struct {
	Node ast.Node
	Ans  map[string]Annotation
}

type Node struct {
	Node ast.Node
}

type Annotation struct {
	Key   string
	Value string
}

type FieldAnnotationVisitor interface {
	Visit(field *ast.Field, ans Annotations)
}

type FuncDeclAnnotationVisitor interface {
	Visit(fn *ast.FuncDecl, ans Annotations)
}

type FileAnnotationVisitor interface {
	Visit(f *ast.File, ans Annotations)
}

type TypeAnnotationVisitor interface {
	Visit(s *ast.TypeSpec, ans Annotations)
}

type AnnotationVisitor struct {
	File     FileAnnotationVisitor
	FuncDecl FuncDeclAnnotationVisitor
	Field    FieldAnnotationVisitor
	Type     TypeAnnotationVisitor
}

func (a *AnnotationVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.File:
		if n.Doc != nil {
			a.File.Visit(n, newAnnotations(n, n.Doc))
		}
	case *ast.Field:
		if n.Doc != nil {
			a.Field.Visit(n, newAnnotations(n, n.Doc))
		}
	case *ast.FuncDecl:
		if n.Doc != nil {
			a.FuncDecl.Visit(n, newAnnotations(n, n.Doc))
		}
	case *ast.TypeSpec:
		if n.Doc != nil {
			a.Type.Visit(n, newAnnotations(n, n.Doc))
		}
	}
	return a
}

func newAnnotations(n ast.Node, cg *ast.CommentGroup) Annotations {
	ans := make(map[string]Annotation)
	for _, c := range cg.List {
		text, ok := extractContent(c)
		if !ok {
			continue
		}
		if strings.HasPrefix(text, "@") {
			segs := strings.SplitN(text, " ", 2)
			if len(segs) != 2 {
				continue
			}
			key := segs[0][1:]
			ans[key] = Annotation{
				Key:   key,
				Value: segs[1],
			}
		}
	}
	return Annotations{
		Node: n,
		Ans:  ans,
	}
}

func extractContent(c *ast.Comment) (string, bool) {
	text := c.Text
	if strings.HasPrefix(text, "// ") {
		return text[3:], true
	} else if strings.HasPrefix(text, "/* ") {
		length := len(text)
		return text[3 : length-2], true
	}
	return "", false
}
