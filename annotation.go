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

type annotations struct {
	Node ast.Node
	Ans  []Annotation
}

type Annotation struct {
	Key   string
	Value string
}

func newAnnotations(n ast.Node, cg *ast.CommentGroup) annotations {
	if cg == nil || len(cg.List) == 0 {
		return annotations{Node: n}
	}
	ans := make([]Annotation, 0, len(cg.List))
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
			ans = append(ans, Annotation{
				Key:   key,
				Value: segs[1],
			})
		}
	}
	return annotations{
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
