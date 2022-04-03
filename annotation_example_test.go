// annotation go through the source code and extra the annotation
// @author Deng Ming
/* @multiple first line
second line
*/
// @date 2022/04/02
package annotation

type (
	// FuncType is a type
	// @author Deng Ming
	/* @multiple first line
	   second line
	*/
	// @date 2022/04/02
	FuncType func()
)

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

// Interface is a test interface
// @author Deng Ming
/* @multiple first line
second line
*/
// @date 2022/04/02
type Interface interface {
	// MyFunc is a test func
	// @parameter arg1 int
	// @parameter arg2 int32
	// @return string
	MyFunc(arg1 int, arg2 int32) string
}
