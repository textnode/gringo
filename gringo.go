// Copyright 2012 Darren Elwood <darren@textnode.com> http://www.textnode.com @textnode
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at 
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Version 0.1
//
// A minimalist queue implemented using a stripped-down lock-free ringbuffer. 
// Inspired by a talk by @mjpt777 at @devbash titled "Lessons in Clean Fast Code" which, 
// among other things, described the LMAX Disruptor.
//
// N.B. To see the performance benefits of gringo versus Go's channels, you must have multiple goroutines
// and GOMAXPROCS > 1.

// Known Limitations:
//
// *) At most (2^64)-1 items can be written to the queue.
// *) The size of the queue must be a power of 2.
// 
// Suggestions:
//
// *) If you have enough cores you can change from runtime.Gosched() to a busy loop.
// *) If you don't need copy-on-write/read semantics change gringo to store pointers.
//
// Warning:
//
// During testing across multiple cores and multiple goroutines the algorithm unexpectedly
// performed correctly WITHOUT the use of additional memory barriers (x86: lfence, mfence, sfence).
//
// If you choose to use this algorithm, test it extensively in your real deployment configurations!
// 

package gringo

import "fmt"
import "sync/atomic"
import "runtime"

// Example item which we will be writing to and reading from the queue
type Payload struct {
    value uint64
    instrument uint64
    price uint64
    quantity uint32
    side bool
    freetext [64]byte
}

func NewPayload(value uint64) *Payload {
    return &Payload{value : value}
}

func (self *Payload) Value () uint64 {
    return self.value
}

// The queue
const queueSize uint64 = 4096
const indexMask uint64 = queueSize - 1

type Gringo struct {
    padding1 [8]uint64
    lastCommittedIndex uint64
    padding2 [8]uint64
    nextFreeIndex uint64
    padding3 [8]uint64
    readerIndex uint64
    padding4 [8]uint64
    contents [queueSize]Payload
    padding5 [8]uint64
}

func NewGringo() *Gringo {
    return &Gringo{lastCommittedIndex : 0, nextFreeIndex : 1, readerIndex : 1}
}

func (self *Gringo) Write(value Payload) {
    var myIndex = atomic.AddUint64(&self.nextFreeIndex, 1) - 1
    //Wait for reader to catch up, so we don't clobber a slot which it is (or will be) reading
    for myIndex > (self.readerIndex + queueSize - 2) {
        runtime.Gosched()
    }
    //Write the item into it's slot
    self.contents[myIndex & indexMask] = value
    //Increment the lastCommittedIndex so the item is available for reading
    for !atomic.CompareAndSwapUint64(&self.lastCommittedIndex, myIndex - 1, myIndex) {
        runtime.Gosched()
    }
}

func (self *Gringo) Read() Payload {
    var myIndex = atomic.AddUint64(&self.readerIndex, 1) - 1
    //If reader has out-run writer, wait for a value to be committed
    for myIndex > self.lastCommittedIndex {
        runtime.Gosched()
    }
    return self.contents[myIndex & indexMask]
}

func (self *Gringo) Dump() {
    fmt.Printf("lastCommitted: %3d, nextFree: %3d, readerIndex: %3d, content:", self.lastCommittedIndex, self.nextFreeIndex, self.readerIndex)
    for index, value := range self.contents {
        fmt.Printf("%5v : %5v", index, value)
    }
    fmt.Print("\n")
}

