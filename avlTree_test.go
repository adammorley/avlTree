package avlTree

import "math"
import "testing"

func TestInsert(test *testing.T) {
    var c int = 5000
    var nums []int = makeNumbers(c)
    createBigTree(nums, c, test)
}

func TestEdges(test *testing.T) {
    /*
        edge case testing (cherry pick from test cases)
    */
}

func TestMany(test *testing.T) {
    /*
        test_many_loop
        test_many_internal
        test_many

    */
}

func makeNumbers(c int) []int {
    var nums []int = make([]int, c+2)
    for i := 0; i < c; i++ {
        nums[i] = i;
    }
    nums[c-1] = math.MaxInt64
    nums[c] = math.MinInt64
    return nums
}

func createBigTree(nums []int, c int, t *testing.T) *avlTree {
    var tree *avlTree = New()
    for i := 0; i <= c; i++ {
        tree.Insert(nums[i])
        if tree.Size() != uint64(i + 1) {
            t.Error("wrong size tree at ", i)
        }
    }
    for i := 0; i <= c; i++ {
        v, e := tree.Search(nums[i])
        if v != nums[i] {
            t.Error("missing number ", v)
        } else if e != nil {
            t.Error("got error ", e)
        }
    }
    return tree
}
