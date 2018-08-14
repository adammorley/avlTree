// implements an AVL tree, a balanced binary search tree.  this implementation stores integers
package avlTree

import "fmt"
import "math"

// a node of the tree.  note that inserting the same value multiple times results in a increment to the counter, not multiple stores
// a node consists of the count of the number of times a value has been stored, a balance factor (to do tree balancing, the pointers to the parent, left and right children, and the value itself.
type node struct {
    /*
        packing into uint16 to fit both into 32bits
    */
    // number of times this interface has been put in the node
    count uint8
    // the balance factor for the node
    bf int8
    parent *node
    left *node
    right *node
    value int
}

type avlTree struct {
    root *node
    // number of nodes; not total count (need to walk tree and sum count)
    size uint64
}

//create a new binary tree
func New() *avlTree {
    var t *avlTree = new(avlTree)
    t.root = nil
    t.size = 0
    return t
}

// create a new node
func newNode(i int) *node {
    var n *node = new(node)
    n.left = nil
    n.right = nil
    n.parent = nil
    n.value = i
    n.count = 1
    n.bf = 0
    return n
}

//insert the provided interface into the tree
func (t *avlTree) Insert(i int) {
    var r *node = t.getRoot()
    if r == nil {
        r = newNode(i)
        t.root = r
        assert(t.size < math.MaxUint64, "too many nodes in tree")
        t.size += 1
        return
    }
    var n *node = t.search(i)
    if n != nil { // value already inserted
        assert(n.count < math.MaxUint8, fmt.Sprintf("too many inserts of value %v", i))
        n.count += 1
        return
    }
    n = newNode(i)
    r = r.insert(n)
    // if the retracing in the insertion case rebalances the tree and changes the root, then update the root pointer in the tree.
    if r.parent == nil {
        t.root = r
    }
    t.size += 1
}

// insert a new node c into the tree starting at node n, maintaining tree balance
func (n *node) insert(c *node) *node {
    if c.value < n.value && n.left != nil {
        return n.left.insert(c)
    } else if c.value > n.value && n.right != nil {
        return n.right.insert(c)
    } else if c.value < n.value && n.left == nil {
        n.left = c
    } else if c.value > n.value && n.right == nil {
        n.right = c
    }
    c.parent = n
    return c.retraceInsert()
}

// node c was just inserted, hence balance is zero; retrace up the tree and rebalance if needed
func (c *node) retraceInsert() *node {
    var p *node = c
    /*
        note break cases:
            if rebalance led to root node
            if tree is balanced (no height change)
    */
    for {
        p = p.parent
        p.updateBalanceFactorInsert(c)
        if p.bf == 2 || p.bf == -2 {
            p = p.rebalance()
        }
        if p.parent == nil {
            break
        }
        if p.bf == 0 {
            break
        }
        c = p
    }
    return p.upToRoot()
}

// update the balance factor of the parent for node n
func (p *node) updateBalanceFactorInsert(c *node) {
    assert(p != nil, "parent is null")
    if p.left == c {
        p.bf -= 1
    } else if p.right == c {
        p.bf += 1
    } else {
        assert(p.left == nil && p.right == nil, "unhandled parent/child relationship")
    }
    assert(p.bf < 3 && p.bf > -3, "balance factor invariant")
}

// rebalance the tree as it is unbalanced, returning the new top node at that level
func (n *node) rebalance() *node {
    if n.bf == 2 { // right heavy
        if n.right.bf == 0 || n.right.bf == 1 { // right right
            n = n.right_right()
        } else { // right left, n.right.bf == -1
            n = n.right_left()
        }
    } else if n.bf == -2 { // left heavy
        if n.left.bf == -1 || n.left.bf == 0 { // left left
            n = n.left_left()
        } else { // left right n.left.bf == 1
            n = n.left_right()
        }
    } else {
        assert(false, "asked to rebalance but bf != 2")
    }
    if n.parent == nil {
        return n
    } else if n.parent.value < n.value {
        n.parent.right = n
    } else if n.parent.value > n.value {
        n.parent.left = n
    } else {
        assert(false, "unhandled parent/child case")
    }
    return n
}

// XXX merge two trees
/*func (t0 *avlTree) Merge(t1 *avlTree) {
}*/

// remove from a tree; if inserted > 1 time, decrement the count
/*func (t *avlTree) Remove(i int) bool {
}*/

type SearchError struct {
    val int
}
func (e *SearchError) Error() string {
    return fmt.Sprintf("could not find value %v in tree", e.val)
}

// search for a value
func (t *avlTree) Search(i int) (int, error) {
    r := t.search(i)
    if r == nil {
        return 0, &SearchError{val: i}
    }
    return r.value, nil
}

// search a tree for a value, return the node with the value
func (t *avlTree) search(i int) *node {
    var n *node = t.getRoot()
    if n == nil {
        return nil
    }
    return n.search(i)
}

// recursive search for a value at a given node
func (n *node) search(i int) *node {
    if i == n.value {
        return n
    } else if i < n.value && n.left != nil {
        return n.left.search(i)
    } else if i > n.value && n.right != nil {
        return n.right.search(i)
    } else {
        return nil
    }
}

func (t *avlTree) Size() uint64 {
    return t.size
}

/*
    tree rebalancing

    there are four unbalanced states:
        right right (the right side is right heavy)
        right left (...)
        left left
        left right

        each of the cases returns the new top node, so the caller can check if this is the new tree root
*/

/*
    right right

        X bf=2
       / \
      t0  Z bf=[0,1]
         / \
        t1  t2
    pivot to:
          Z bf=[-1,0]
         / \
bf=[0,1]X   t2
       / \
      t0 t1
*/
func (X *node) right_right() *node {
    var Z *node = X.right
    assert(Z.bf != -1, "right left in right right case")
    Z.parent = X.parent
    X.parent = Z

    X.right = Z.left
    if X.right != nil {
        X.right.parent = X
    }

    Z.left = X
    if Z.bf == 0 { // delete from X left
        X.bf = 1
        Z.bf = -1
    } else { // insert R or delete L from Z
        X.bf = 0
        Z.bf = 0
    }
    return Z
}

/*
    right left

          X bf=2
         / \
        t0  Z bf=-1
           / \
bf=-1,0,1 Y  t3
         / \
        t1 t2
    pivot to:
           Y bf=0
          / \
bf=-1,0  X   Z bf=0,1
        / \  / \
      t0 t1 t2 t3
*/
func (X *node) right_left() *node {
    var Z *node = X.right
    assert(Z.bf == -1, "right right in left right case")
    var Y *node = Z.left
    Y.parent = X.parent
    Z.parent = Y
    X.parent = Y

    X.right = Y.left
    if X.right != nil {
        X.right.parent = X
    }
    Z.left = Y.right
    if Z.left != nil {
        Z.left.parent = Z
    }

    Y.left = X
    Y.right = Z
    if Y.bf == 0 { // delete from X L
        X.bf = 0
        Z.bf = 0
    } else if Y.bf == -1 { // insert on Y L
        X.bf = 0
        Z.bf = 1
        Y.bf = 0
    } else if Y.bf == 1 { // insert on Y R
        X.bf = -1
        Z.bf = 0
        Y.bf = 0
    } else {
        assert(false, "invalid balance factor for Y in right left case")
    }
    return Y
}

/*
    left left

           X bf=-2
          / \
bf=[-1,0] Z   t2
        / \
       t0 t1
    pivot to:
        Z bf=[0,1]
       / \
      t0  X bf=[-1,0]
         / \
        t1  t2
*/
func (X *node) left_left() *node {
    var Z *node = X.left
    assert(Z.bf != 1, "left right case in left left")
    Z.parent = X.parent
    X.parent = Z

    X.left = Z.right
    if X.left != nil {
        X.left.parent = X
    }

    Z.right = X
    if Z.bf == 0 { // delete from X R
        X.bf = -1
        Z.bf = 1
    } else { // insert L or delete R from Z
        X.bf = 0
        Z.bf = 0
    }
    return Z
}

/*
    left right

        X bf=-2
       / \
 bf=1 Z   t3
     / \
    t0  Y bf=[-1,0,1]
       / \
      t1 t2
    pivot to:
        Y bf=0
       / \
bf=0,1 Z   X bf=-1,0
     / \ / \
   t0 t1 t2 t3
*/
func (X *node) left_right() *node {
    var Z *node = X.left
    assert(Z.bf == 1, "left left case in left right")
    var Y *node = Z.right
    Y.parent = X.parent
    Z.parent = Y
    X.parent = Y

    X.left = Y.right
    if X.left != nil {
        X.left.parent = X
    }
    Z.right = Y.left
    if Z.right != nil {
        Z.right.parent = Z
    }

    Y.right = X
    Y.left = Z
    if Y.bf == 0 { // delete from X R
        X.bf = 0
        Z.bf = 0
    } else if Y.bf == -1 { // insert on Y L
        X.bf = 0
        Z.bf = 1
        Y.bf = 0
    } else if Y.bf == 1 { // insert on Y R
        X.bf = -1
        Z.bf = 0
        Y.bf = 0
    } else {
        assert(false, "unhandled balance factor for Y in left right case")
    }
    return Y
}

// support assert calls
func assert(condition bool, msg string) {
    if !condition {
        panic(msg)
    }
}

// get the root node of a tree
func (t *avlTree) getRoot() *node {
    assert(t != nil, "tree is null")
    return t.root
}

// iterate up to the root of the tree
func (n *node) upToRoot() *node {
    if n == nil {
        return nil
    }
    for n.parent != nil {
        n = n.parent
    }
    return n
}
