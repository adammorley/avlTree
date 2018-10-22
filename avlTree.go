// implements an AVL tree, a balanced binary search tree for any arbitrary comparable interface
// XXX could add: pre-order, post-order, tree merge
package avlTree

import (
	"errors"
	"fmt"
	"log"
	"math"
)

/*
    consumers of this package must implement the following interfaces in order for the AVL tree to order the elements of the tree.  it's critical to use a type assertion; as a type assertion can cause a panic if the types don't match, optionally one can use a type-testing assertion (eg: i, ok := j.(TYPE) if !ok ...).  however, this would require a minor modification to the Interface of this package, as it would need to allow for error handling (eg: type mismatch)

type MyInt int

func (i MyInt) LessThan(j interface{}) bool {
    return i < j.(MyInt)
}

func (i MyInt) GreaterThan(j interface{}) bool {
    return i > j.(MyInt)
}

func (i MyInt) EqualTo(j interface{}) bool {
    return i == j.(MyInt)
}
*/

type Interface interface {
	LessThan(j interface{}) bool
	GreaterThan(j interface{}) bool
	EqualTo(j interface{}) bool
}

// a node of the tree.  note that inserting the same value multiple times results in a increment to the counter, not multiple stores
// a node consists of the count of the number of times a value has been stored, a balance factor (to do tree balancing), the pointers to the parent, left and right children, and the value itself.
type node struct {
	//packing into uint8 to fit both into 16bits
	count  uint8 // number of times this interface has been put in the node
	bf     int8  // the balance factor for the node
	parent *node
	left   *node
	right  *node
	value  Interface
}

const sizeBits = 64 // avoid using unsafe
type avlTree struct {
	root *node
	size uint64 // number of nodes; not total count (need to walk tree and sum count)
}

//create a new binary tree
func New() (t *avlTree) {
	t = new(avlTree)
	t.root = nil
	t.size = uint64(0)
	return
}

// create a new node
func newNode(i Interface) (n *node) {
	n = new(node)
	n.left = nil
	n.right = nil
	n.parent = nil
	n.value = i
	n.count = 1
	n.bf = 0
	return
}

/*
   NODE VISITATION IN ORDER
*/

type stack struct {
	pile []*node
	top  uint
}

// get the height of the tree by finding the highest set bit in size and moving it over one.  this accounts for the dangling node case (eg a height 3 tree is 5 nodes, eg 5 (101) -> 8 (1000) -> log_2(8) == 3
func (t *avlTree) getHeight() uint {
	if t.size == 0 {
		return 0
	}
	var i uint64
	for i = sizeBits - 1; i >= 0; i-- {
		set := (1 << i) & t.size
		if set>>i == 1 {
			twoToTheH := 1 << (i + 1)
			return uint(math.Log2(float64(twoToTheH)))
		}
	}
	assert(false, "cannot reach here")
	return 0
}

// Inorder returns a function which will return the next node in order in the tree.  This is done iteratively using a stack and pointer model
func (t *avlTree) Inorder() func() (Interface, error) {
	s := new(stack)
	// the maximum size of the stack is the height of the tree
	s.pile = make([]*node, t.getHeight())
	n := t.getRoot()
	return func() (Interface, error) {
		for s.top > 0 || n != nil {
			if n != nil {
				s.pile[s.top] = n
				s.top = s.top + 1
				n = n.left
			} else {
				s.top = s.top - 1
				n = s.pile[s.top].right
				return s.pile[s.top].value, nil
			}
		}
		return nil, errors.New("end of tree")
	}
}

/*
   NODE INSERTION
*/

//insert the interface into the tree
func (t *avlTree) Insert(i Interface) {
	assert(t.size < math.MaxUint64, "too many nodes in tree")
	if t.getRoot() == nil { // root case
		assert(t.size == 0, "wrong size")
		t.root = newNode(i)
		t.size += 1
		return
	}

	var n *node = t.searchForClosest(i) // won't return nil
	if n.value.EqualTo(i) {             // value already inserted
		assert(n.count < math.MaxUint8, fmt.Sprintf("too many inserts of value %v", i))
		n.count += 1
		return
	}

	c := newNode(i) // new node case
	if c.value.LessThan(n.value) && n.left == nil {
		n.left = c
	} else if c.value.GreaterThan(n.value) && n.right == nil {
		n.right = c
	} else {
		log.Fatal("invariant", c, n, n.left, n.right)
	}
	c.parent = n
	n = retraceInsert(c)

	// if the retracing in the insertion case rebalances the tree and changes the root, then update the root pointer in the tree.
	if n.parent == nil {
		t.root = n
	}
	t.size += 1
}

// node c was just inserted, hence balance is zero; retrace up the tree and rebalance if needed
func retraceInsert(c *node) *node {
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
			p = rebalance(p)
		}
		if p.parent == nil {
			break
		}
		if p.bf == 0 {
			break
		}
		c = p
	}
	return upToRoot(p)
}

// update the balance factor of the parent for node n
func (p *node) updateBalanceFactorInsert(c *node) {
	assert(p != nil, "parent is nil")
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
func rebalance(n *node) *node {
	if n.bf == 2 { // right heavy
		if n.right.bf == 0 || n.right.bf == 1 { // right right
			n = right_right(n)
		} else { // right left, n.right.bf == -1
			n = right_left(n)
		}
	} else if n.bf == -2 { // left heavy
		if n.left.bf == -1 || n.left.bf == 0 { // left left
			n = left_left(n)
		} else { // left right n.left.bf == 1
			n = left_right(n)
		}
	} else {
		assert(false, "asked to rebalance but bf != 2")
	}
	if n.parent == nil {
		return n
	} else if n.parent.value.LessThan(n.value) {
		n.parent.right = n
	} else if n.parent.value.GreaterThan(n.value) {
		n.parent.left = n
	} else {
		assert(false, "unhandled parent/child case")
	}
	return n
}

/*
   VALUE DELETION
*/

// remove from a tree; if inserted > 1 time, decrement the count instead of removing the node
func (t *avlTree) Delete(i Interface) bool {
	var r *node = t.getRoot()

	// handle root case
	if r == nil {
		return false
	} else if r.value.EqualTo(i) && r.left == nil && r.right == nil && r.count == 1 {
		assert(t.Size() == 1, "root deletion, but size wrong")
		t.root = nil
		t.size = 0
		return true
	}

	// handle value stored more than once case
	var n *node = t.search(i)
	if n == nil {
		return false
	} else if n.count > 1 {
		n.count -= 1
		return true
	}

	// handle node removal (last stored value)
	if n.right == nil && n.left == nil {
		r = removeNoChildren(n)
	} else if n.right == nil {
		r = removeNoRightChildren(n)
	} else if n.right != nil && n.right.left == nil {
		r = removeRightNoLeft(n)
	} else if n.right != nil && n.right.left != nil {
		r = removeComplex(n)
	} else {
		assert(false, "unhandled node removal")
	}

	// update root
	if r == nil {
		t.root = nil
	} else if r.parent != nil {
		assert(false, "root parent not nil")
	}
	t.root = r
	t.size -= 1
	return true
}

// replace node c with x at p
func replaceNode(c, p, x *node) {
	if x != nil {
		x.parent = p
	}
	if p == nil { // root case
		return
	} else if p.left == c {
		p.left = x
	} else if p.right == c {
		p.right = x
	} else {
		assert(false, "unhandled parent/child case")
	}
}

// remove a node with no children, simply update the balance factors and retrace
func removeNoChildren(n *node) *node {
	var p *node = n.parent
	if p == nil { // root node
		return nil
	} else if p.right == n {
		p.bf -= 1
	} else {
		p.bf += 1
	}
	replaceNode(n, p, nil)
	return retraceRemove(p)
}

/*
    remove a node with no right children, simply swap children

    10
   /  \
  3   20*
     /
    15

    to

    10
   /  \
  3   15
*/
func removeNoRightChildren(n *node) *node {
	var p, x *node = n.parent, n.left
	replaceNode(n, p, x)
	if p == nil && x.parent == nil { // root node
		return x
	} else {
		return retraceRemove(x)
	}
}

/*
    remove node with right child, but child has no left children

    10
   /  \
  3   20
     /  \
    15  25*
          \
          30
            \
            40

    to

    10
   /  \
  3   20
     /  \
    15   30
           \
           40
*/
func removeRightNoLeft(n *node) *node {
	var p, x *node = n.parent, n.right
	x.left = n.left
	x.bf = n.bf - 1
	replaceNode(n, p, x)
	if p == nil && x.parent == nil { // root node
		return x
	} else {
		return retraceRemove(x)
	}
}

/*
    node has a right child and child has left children; traverse down and re-home left most child.  don't drop right child of left most child!

    10
   /  \
  3   20
     /  \
    15  25*
       /  \
      24  30
         /  \
        27  40
         \
         28

    to

    10
   /  \
  3   20
     /  \
    15  27
       /  \
      24  30
         /  \
        28  40
*/
func removeComplex(n *node) *node {
	var c *node = n.right
	// find the availble left node
	for c.left != nil {
		c = c.left
	}

	// find the child
	if c.right != nil {
		c.parent.left = c.right
	} else {
		c.parent.left = nil
	}

	// height shrinks
	var p *node = c.parent
	p.bf += 1

	// put the child in the right place
	c.bf = n.bf
	replaceNode(n, n.parent, c)

	// swap in the children from n to c
	n.left.parent = c
	c.left = n.left
	n.right.parent = c
	c.right = n.right

	// retrace
	if p == nil && c.parent == nil {
		return c
	} else {
		return retraceRemove(p)
	}
}

// node n just had a child removed, retrace & rebalance if needed, returning new root node
func retraceRemove(n *node) *node {
	/*
	   note the break cases:
	       if node removal was absorbed at n (combined with case below)
	       if after rebalance, removal was absorbed
	       if node doesn't have parent (root node)
	       otherwise, update balance factors, check asserts, loop
	*/
	for {
		if n.bf == 2 || n.bf == -2 {
			n = rebalance(n)
		}
		assert(n.bf < 2 && n.bf > -2, "node balance factor out of range after rebalance")
		if n.bf == 1 || n.bf == -1 {
			break
		} else if n.parent == nil {
			break
		} else if n.parent.value.LessThan(n.value) {
			n.parent.bf -= 1
		} else if n.parent.value.GreaterThan(n.value) {
			n.parent.bf += 1
		} else {
			assert(false, "unhandled parent/child relationship")
		}
		assert(n.parent.bf < 3 && n.parent.bf > -3, "balance factor invariant")
		n = n.parent
	}
	assert(n != nil, "n nil")
	return upToRoot(n)
}

/*
   NODE SEARCH
*/

type SearchError struct {
	val Interface
}

func (e *SearchError) Error() string {
	return fmt.Sprintf("could not find value %v in tree", e.val)
}

// search for a value
func (t *avlTree) Search(i Interface) (Interface, error) {
	var n *node = t.search(i)
	if n == nil {
		return nil, &SearchError{val: i}
	}
	return n.value, nil
}

// search a tree for a value, return the node with the value, or nil
func (t *avlTree) search(i Interface) *node {
	n := t.searchForClosest(i)
	if n == nil || !i.EqualTo(n.value) {
		return nil
	} else {
		return n
	}
}

// search for the closest node for the interface, using comparison (this allows re-use of the code by the insert case)
func (t *avlTree) searchForClosest(i Interface) *node {
	var n *node = t.getRoot()
	if n == nil {
		return nil
	}
	for !i.EqualTo(n.value) && (i.LessThan(n.value) && n.left != nil || i.GreaterThan(n.value) && n.right != nil) {
		if i.LessThan(n.value) {
			n = n.left
		} else if i.GreaterThan(n.value) {
			n = n.right
		} else {
			log.Fatal("invariant", i, n, n.left, n.right)
		}
	}
	return n
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
func right_right(X *node) (Z *node) {
	Z = X.right
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
	return
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
func right_left(X *node) (Y *node) {
	var Z *node = X.right
	assert(Z.bf == -1, "right right in left right case")
	Y = Z.left
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
	return
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
func left_left(X *node) (Z *node) {
	Z = X.left
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
	return
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
func left_right(X *node) (Y *node) {
	var Z *node = X.left
	assert(Z.bf == 1, "left left case in left right")
	Y = Z.right
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
	return
}

// support assert calls
func assert(condition bool, msg string) {
	if !condition {
		log.Fatal(msg)
	}
}

// get the root node of a tree
func (t *avlTree) getRoot() *node {
	assert(t != nil, "tree is nil")
	return t.root
}

// iterate up to the root of the tree
func upToRoot(n *node) *node {
	for n.parent != nil {
		n = n.parent
	}
	return n
}
