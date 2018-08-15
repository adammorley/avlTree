package avlTree

import "container/list"
import "math"
import "testing"

type Int int
var MaxInt Int = Int(math.MaxInt64)
var MinInt Int = Int(math.MinInt64)

func (i Int) LessThan(j interface{}) bool {
    return i < j.(Int)
}

func (i Int) GreaterThan(j interface{}) bool {
    return i > j.(Int)
}

func (i Int) EqualTo(j interface{}) bool {
    return i == j.(Int)
}

func TestDoubleInsertAndDelete(test *testing.T) {
	var tree *avlTree = New()

	insertOneAndCheck(test, tree)
	insertOneAndCheck(test, tree)

    var one Int = Int(1)
	if !tree.Delete(one) {
		test.Error("deletion does not work")
	} else if _, e := tree.Search(one); e != nil {
		test.Error("couldn't find value")
	} else if tree.Size() != 1 {
		test.Error("wrong size")
	}

	if !tree.Delete(one) {
		test.Error("deletion does not work")
	} else if tree.Size() != 0 {
		test.Error("wrong size")
	}
	_, e := tree.Search(one)
	if e == nil {
		test.Error("should not find value")
	} else if _, ok := e.(*SearchError); !ok {
		test.Error("type returned from search incorrect", e)
	}
}

func insertOneAndCheck(test *testing.T, tree *avlTree) {
    var one Int = Int(1)
	tree.Insert(one)
	if _, e := tree.Search(one); e != nil {
		test.Error("couldn't find value")
	} else if tree.Size() != 1 {
		test.Error("wrong size")
	}
}

func TestBigTree(test *testing.T) {
	var c int = 5000
	var nums []Int = makeNumbers(c)
	var tree *avlTree = createBigTree(nums, c, test)
	checkBalanceFactor(tree.getRoot(), test)
	for i := 0; i <= c; i++ {
		v, e := tree.Search(nums[i])
		if v != nums[i] {
			test.Error("missing number ", v)
		} else if e != nil {
			test.Error("got error ", e)
		}
	}

}

func makeNumbers(c int) []Int {
	var nums []Int = make([]Int, c+2)
	for i := 0; i < c; i++ {
		nums[i] = Int(i)
	}
	nums[c-1] = MaxInt
	nums[c] = MinInt
	return nums
}

func createBigTree(nums []Int, c int, test *testing.T) *avlTree {
	var tree *avlTree = New()
	for i := 0; i <= c; i++ {
		tree.Insert(nums[i])
		if tree.Size() != uint(i+1) {
			test.Error("wrong size tree at ", i)
		}
	}
	return tree
}

func checkBalanceFactor(n *node, test *testing.T) {
	if n == nil {
		return
	} else if n.bf > 2 || n.bf < -2 {
		test.Error("node ", n.value, " has incorrect balance factor ", n.bf)
	}
	checkBalanceFactor(n.left, test)
	checkBalanceFactor(n.right, test)
}

func TestMany(test *testing.T) {
	for i := 1; i <= 999; i++ {
		testManyInternal(i, test)
	}
}

func testManyLoop(s int, k int, c int, n []Int, test *testing.T) {
	var tree *avlTree = createBigTree(n, c, test)
	for i := s; i <= c-2; i += 2 {
        var x Int = Int(i)
		if !tree.Delete(x) {
			test.Error("could not remove", i)
		}
	}
	checkBalanceFactor(tree.getRoot(), test)
	for i := k; i <= c-2; i += 2 {
        var x Int = Int(i)
		if _, e := tree.Search(x); e != nil {
			test.Error("could not find", x)
		}
	}
	if _, e := tree.Search(MaxInt); e != nil {
		test.Error("could not find", MaxInt)
	} else if _, e := tree.Search(MinInt); e != nil {
		test.Error("could not find", MinInt)
	}
}

func testManyInternal(c int, test *testing.T) {
	var n []Int = makeNumbers(c)
	testManyLoop(0, 1, c, n, test)
	testManyLoop(1, 0, c, n, test)
}

func TestInorder(test *testing.T) {
	var tree *avlTree = New()
	var s, f int = 0, 10
	for i := s; i <= f; i++ {
        var x Int = Int(i)
		tree.Insert(x)
	}
	var l *list.List = tree.Inorder()
	var e *list.Element = l.Front()
	for i := s; i <= f; i++ {
        var x Int = Int(i)
		if e.Value != x {
			test.Error("wrong value", e.Value, x)
		}
		e = e.Next()
	}
}
