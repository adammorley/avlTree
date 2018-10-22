package avlTree

import (
	"math"
	"testing"
)

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
	_, err := tree.Search(one)
	if err == nil {
		test.Error("should not find value")
	} else if _, ok := err.(*SearchError); !ok {
		test.Error("type returned from search incorrect", err)
	}
}
func insertOneAndCheck(test *testing.T, tree *avlTree) {
	var one Int = Int(1)
	tree.Insert(one)
	if _, err := tree.Search(one); err != nil {
		test.Error("couldn't find value")
	} else if tree.Size() != 1 {
		test.Error("wrong size")
	}
}

func TestHeight(test *testing.T) {
	tree := new(avlTree)
	tree.Insert(Int(1))
	if tree.getHeight() != 1 {
		test.Error("wrong height")
	}
	tree.Insert(Int(2))
	if tree.getHeight() != 2 {
		test.Error("wrong height")
	}
	tree.Insert(Int(3))
	tree.Insert(Int(4))
	if tree.getHeight() != 3 {
		test.Error("wrong height")
	}
}

func TestInorder(test *testing.T) {
	tree := New()
	var s, f int = 0, 10
	for i := s; i <= f; i++ {
		tree.Insert(Int(i))
	}
	inorder := tree.Inorder()
	for i := s; i <= f; i++ {
		v, err := inorder()
		if err != nil {
			test.Error("should've gotten no error", err)
		} else if v != Int(i) {
			test.Error("wrong value", v, Int(i))
		}
	}
	_, err := inorder()
	if err == nil {
		test.Error("should've reached end of tree")
	}
}

func TestBigTree(test *testing.T) {
	c := 5000
	var nums []Int = makeNumbers(c)
	var tree *avlTree = createBigTree(nums, c, test)
	checkBalanceFactor(tree.getRoot(), test)
	for i := 0; i <= c; i++ {
		v, err := tree.Search(nums[i])
		if v != nums[i] {
			test.Error("missing number ", v)
		} else if err != nil {
			test.Error("got error ", err)
		}
	}

}
func makeNumbers(c int) (nums []Int) {
	nums = make([]Int, c+2)
	for i := 0; i < c; i++ {
		nums[i] = Int(i)
	}
	nums[c-1] = MaxInt
	nums[c] = MinInt
	return
}
func createBigTree(nums []Int, c int, test *testing.T) (tree *avlTree) {
	tree = New()
	for i := 0; i <= c; i++ {
		tree.Insert(nums[i])
		if tree.Size() != uint64(i+1) {
			test.Error("wrong size tree at ", i)
		}
	}
	return
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

/*
   test out many different tree cases
*/
func TestMany(test *testing.T) {
	for i := 1; i <= 999; i++ {
		testManyInternal(i, test)
	}
}
func testManyInternal(c int, test *testing.T) {
	var n []Int = makeNumbers(c)
	for i := 2; i < 8; i++ {
		testManyLoop(0, 1, c, i, n, test)
		testManyLoop(1, 0, c, i, n, test)
	}
}

/*
   loop through the tree and delete values, then ensure nothing else
   got dropped during deletion
*/
func testManyLoop(s, k, c, p int, n []Int, test *testing.T) {
	var tree *avlTree = createBigTree(n, c, test)
	for i := s; i <= c-2; i += p {
		if !tree.Delete(Int(i)) {
			test.Error("could not remove", i)
		}
	}
	checkBalanceFactor(tree.getRoot(), test)
	for i := k; i <= c-2; i += p {
		if _, err := tree.Search(Int(i)); err != nil {
			test.Error("could not find", i)
		}
	}
	if _, err := tree.Search(MaxInt); err != nil {
		test.Error("could not find", MaxInt)
	} else if _, err := tree.Search(MinInt); err != nil {
		test.Error("could not find", MinInt)
	}
}
