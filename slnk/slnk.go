// package slnk - singly linked list
// Supports:
// - FIFO and LIFO opperation
// - variable data element types allowed
// - Access to middle of list by index
// - Deletion from middle of list by index
// - Range limits by index for read and print access
// - Printing with optional Debug option
package slnk

import (
	"errors"
	"fmt"
	"io"
	"unsafe"
)

// singally linked list type
// represents the anchor point and all common statust information for the list
type slist struct {
	head  *node                              // size=8  addr first node in list
	tail  *node                              // size=8  addr last node of list
	curr  *node                              // size=8  addr current node after locate
	cnt   int                                // size=8  number of nodes
	cmp   func(interface{}, interface{}) int // size=8  compare func, ret -1,0,1 : lt,eq,gt
	err   error                              // size=16 last error for method or nil
	debug bool                               // size=1  debug enabled flag
	start int                                // size=8  start index of range parameter for some methods
	end   int                                // size=8  end index of range parameter for some methods
}

// linked list element.
// because the data type for a list data element is interface{} (aka a variant)
// each element in the list can be a different type
type node struct {
	next *node       // size=8  address of next node
	data interface{} // size=16 reference tonode data payload
}

// create a new singally linked list
func NewSlist() *slist {
	s := new(slist)
	return s
}

// Size returns the current size in bytes of the list
func (s *slist) Size() int {
	// size of header
	v := slist{}
	size := unsafe.Sizeof(v)
	// get size of nodes and data
	d := node{}
	for n := s.head; n != nil; n = n.next {
		size += unsafe.Sizeof(d)
		size += unsafe.Sizeof(n.data)
	}
	return int(size)
}

// setCount sets the count of items in the list and resets range start and end
func (s *slist) setCount(count int) {
	if count < 0 {
		count = 0
	}
	s.cnt = count
	s.start = 0
	if count > 0 {
		s.end = count - 1
	} else {
		s.end = 0
	}
}

// incrCount increments the count of items in the list and resets range start and end
func (s *slist) incrCount() {
	s.cnt++
	s.start = 0
	s.end = s.cnt - 1
}

// decrCount decrements the count of items in the list and resets range start and end
func (s *slist) decrCount() {
	s.start = 0
	if s.cnt <= 1 {
		s.cnt = 0
		s.end = 0
		return
	}
	s.cnt--
	s.end = s.cnt - 1
}

// Error returns the last error from the previous slist function or nil if
// no error occurred
func (s *slist) Error() error {
	return s.err
}

// SetCompareFunc - set the compare function for this slist to use for
// ordered insertion
func (s *slist) SetCompareFunc(f func(interface{}, interface{}) int) *slist {
	s.cmp = f
	return s
}

// Len returns the current length of the list
func (s *slist) Len() int {
	return s.cnt
}

// Add an element to the tail of the linked list
func (s *slist) Add(data interface{}) *slist {
	s.err = nil
	// create new node
	n := new(node)
	n.data = data
	// add new node
	if s.head == nil {
		s.head = n
		s.tail = n
		s.setCount(1) // and reset any range parameters
	} else {
		s.head.next = n
		s.tail = n
		s.incrCount() // and reset any range parameters
	}
	return s
}

// Insert an element after the specified index into the linked list
func (s *slist) Insert(data interface{}, index int) *slist {
	s.err = nil
	// create new node
	n := new(node)
	n.data = data
	// insert new first node
	if s.head == nil {
		if index == 0 {
			s.head = n
			s.tail = n
			s.setCount(1) // and reset any range parameters
			return s
		}
		s.setCount(0)
		s.err = errors.New("cant insert at non zero index in empty list")
		return s
	}
	p := 0
	// find location and insert new node
	for m := s.head; m != nil; m = m.next {
		if p == index {
			tmp := m.next
			m.next = n
			n.next = tmp
			if tmp == nil {
				s.tail = n
			}
			s.incrCount() // and reset any range parameters
		}
		p++
	}
	return s
}

// Delete an element after the specified index from the linked list
func (s *slist) Delete(index int) *slist {
	s.err = nil
	// check range
	if index > s.cnt {
		s.err = errors.New("index outside list range")
		return s
	}
	if s.head == nil {
		s.err = errors.New("cant delete from empty list")
		return s
	}
	// find and delete node
	p := 0
	for m := s.head; m != nil; m = m.next {
		if p == index {
			if m.next != nil {
				m.next = m.next.next
				s.decrCount() // and reset any range parameters
				return s
			}
		}
		p++
	}
	return s
}

// locates element index in the list and saves location.
// if not in range error is set, chec with
// the first element is element zero. if the element is not present nil is
// returned
func (s *slist) Locate(index int) *slist {
	s.err = nil
	s.curr = nil
	if index > s.cnt {
		s.err = errors.New("index outside list range")
		return s
	}
	p := 0
	for n := s.head; n != nil; n = n.next {
		if p == index {
			s.curr = n
			return s
		}
		p++
	}
	return s
}

// Data gets the data from the current element or nil if a current element
// has not been set
func (s *slist) Data() interface{} {
	s.err = nil
	if s.curr == nil {
		s.err = errors.New("no current element set")
		return nil
	}
	return s.curr.data
}

// removes an element from tail of the list
func (s *slist) Remove() *slist {
	s.err = nil
	// check range
	if s.head == nil {
		s.err = errors.New("no elements in list")
		s.cnt = 0
		return s
	} else {
		s.Locate(s.cnt - 1) // locate element before tail
		if s.curr != nil {
			n := s.tail
			s.tail = s.curr
			s.curr = n    // current becomes old tail
			s.decrCount() // and reset any range parameters
		}
		return s
	}
}

// pops an element from the front of the list
func (s *slist) Pop() *slist {
	s.err = nil
	// find element and pop
	if s.head != nil {
		s.curr = s.head
		s.head = s.head.next
		s.decrCount() // and reset any range parameters
		return s
	}
	s.err = errors.New("no elements in list")
	s.cnt = 0
	return s
}

// Print the entire slist in a formatted way
// if range is set limit printing nodes to the specified range
// if debug is set extra information is printed
func (s *slist) Print(o io.Writer) *slist {
	fmt.Fprintf(o, "Slist @: %08p Len: %d\n", s, s.cnt)
	if s.debug {
		fmt.Fprintf(o, "  Range start: %d end: %d\n", s.start, s.end)
		fmt.Fprintf(o, "  Compare func: %08p\n", s.cmp)
		fmt.Fprintf(o, "  Last error: %v\n", s.err)
	}
	if s.start != 0 || s.end != s.cnt {
		fmt.Fprintf(o, "  Nodes in range %d:%d ...\n", s.start, s.end)
	} else {
		fmt.Fprintf(o, "  Nodes ...\n")
	}
	i := 0
	for n := s.head; n != nil; n = n.next {
		if i >= s.start && i <= s.end {
			if s.debug {
				fmt.Fprintf(o, "    %06d,%08p: %v\n      next:%08p\n",
					i, n, n.data, n.next)
			} else {
				fmt.Fprintf(o, "    %06d: %v\n", i, n.data)
			}
		}
		i++
	}
	fmt.Fprintf(o, "%d nodes listed\n", i)
	return s
}

// Debug parameter setting
func (s *slist) Debug(on bool) *slist {
	s.debug = on
	return s
}

// Start of range parameter setting
func (s *slist) Start(start int) *slist {
	if start >= 0 && start <= (s.cnt-2) {
		s.start = start
	} else {
		s.start = 0
		s.err = errors.New("start parameter outside list bounds")
	}
	return s
}

// End of range parameter setting
func (s *slist) End(end int) *slist {
	if end >= 1 && end <= (s.cnt-1) {
		if end > s.start {
			s.end = end
		} else {
			s.err = errors.New("end parameter less than start parameter")
		}
	} else {
		s.end = s.cnt - 1
		s.err = errors.New("end parameter outside list bounds")
	}
	return s
}
