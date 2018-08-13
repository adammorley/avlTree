package avlTree

import "testing"

func TestInsert(test *testing.T) {
    var tree *avlTree = New()
    tree.Insert(1)
    tree.Insert(1)
    tree.Insert(2)
    test.Error(tree.Size())
    test.Error(tree.Search(1))
    test.Error(tree.Search(2))
    test.Error(tree.Search(3))
}
