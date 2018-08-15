package avlTree

import "container/list"
import "math"
import "testing"

func TestDoubleInsertAndDelete(test *testing.T) {
	var tree *avlTree = New()

	insertOneAndCheck(test, tree)
	insertOneAndCheck(test, tree)

	if !tree.Delete(1) {
		test.Error("deletion does not work")
	} else if _, e := tree.Search(1); e != nil {
		test.Error("couldn't find value")
	} else if tree.Size() != 1 {
		test.Error("wrong size")
	}

	if !tree.Delete(1) {
		test.Error("deletion does not work")
	} else if tree.Size() != 0 {
		test.Error("wrong size")
	}
	_, e := tree.Search(1)
	if e == nil {
		test.Error("should not find value")
	} else if _, ok := e.(*SearchError); !ok {
		test.Error("type returned from search incorrect", e)
	}
}

func insertOneAndCheck(test *testing.T, tree *avlTree) {
	tree.Insert(1)
	if _, e := tree.Search(1); e != nil {
		test.Error("couldn't find value")
	} else if tree.Size() != 1 {
		test.Error("wrong size")
	}
}

func TestBigTree(test *testing.T) {
	var c int = 5000
	var nums []int = makeNumbers(c)
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

func makeNumbers(c int) []int {
	var nums []int = make([]int, c+2)
	for i := 0; i < c; i++ {
		nums[i] = i
	}
	nums[c-1] = math.MaxInt64
	nums[c] = math.MinInt64
	return nums
}

func createBigTree(nums []int, c int, test *testing.T) *avlTree {
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

func testManyLoop(s int, k int, c int, n []int, test *testing.T) {
	var tree *avlTree = createBigTree(n, c, test)
	for i := s; i <= c-2; i += 2 {
		if !tree.Delete(i) {
			test.Error("could not remove", i)
		}
	}
	checkBalanceFactor(tree.getRoot(), test)
	for i := k; i <= c-2; i += 2 {
		if _, e := tree.Search(i); e != nil {
			test.Error("could not find", i)
		}
	}
	if _, e := tree.Search(math.MaxInt64); e != nil {
		test.Error("could not find", math.MaxInt64)
	} else if _, e := tree.Search(math.MinInt64); e != nil {
		test.Error("could not find", math.MinInt64)
	}
}

func testManyInternal(c int, test *testing.T) {
	var n []int = makeNumbers(c)
	testManyLoop(0, 1, c, n, test)
	testManyLoop(1, 0, c, n, test)
}

func TestInorder(test *testing.T) {
	var tree *avlTree = New()
	var s, f int = 0, 10
	for i := s; i <= f; i++ {
		tree.Insert(i)
	}
	var l *list.List = tree.Inorder()
	var e *list.Element = l.Front()
	for i := s; i <= f; i++ {
		if e.Value != i {
			test.Error("wrong value", e.Value, i)
		}
		e = e.Next()
	}
}
