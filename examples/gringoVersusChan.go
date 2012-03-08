// Copyright 2012 Darren Elwood <darren@textnode.com> http://www.textnode.com @textnode
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "time"
    "fmt"
    "gringo"
)

var size uint64 = 5000000

func noopPayload(param gringo.Payload) {}

var pl gringo.Payload = *gringo.NewPayload(1)

func gringoProducer(outGringo *gringo.Gringo, done chan int) {
    var i uint64
    for ; i < size; i++ {
        outGringo.Write(pl)
    }
    done <- 0
}

func gringoForwarder(inGringo *gringo.Gringo, outGringo *gringo.Gringo, done chan int) {
    var i uint64
    for ; i < size; i++ {
        var rez gringo.Payload = inGringo.Read()
        outGringo.Write(rez)
    }
    done <- 0
}

func gringoConsumer(inGringo *gringo.Gringo, done chan int) {
    var i uint64
    for ; i < size; i++ {
        var rez gringo.Payload = inGringo.Read()
        noopPayload(rez)
    }
    done <- 0
}

func gringoRunner() {
    doneChan := make(chan int)

    var r1 *gringo.Gringo = gringo.NewGringo()
    var r2 *gringo.Gringo = gringo.NewGringo()

    var startTime time.Time = time.Now()

    go gringoProducer(r1, doneChan)
    go gringoForwarder(r1, r2, doneChan)
    go gringoConsumer(r2, doneChan)

    <-doneChan; <-doneChan; <-doneChan

    fmt.Println("gringoRunner Nanoseconds passed:", time.Since(startTime))
}


func chanProducer(outChan chan gringo.Payload, done chan int) {
    var i uint64
    for ; i < size; i++ {
        outChan <- pl
    }
    done <- 0
}

func chanForwarder(inChan chan gringo.Payload, outChan chan gringo.Payload, done chan int) {
    var i uint64
    for ; i < size; i++ {
        o := <- inChan
        outChan <- o
    }
    close(inChan)
    done <- 0
}

func chanConsumer(inChan chan gringo.Payload, done chan int) {
    var i uint64
    for ; i < size; i++ {
        o := <- inChan
        noopPayload(o)
    }
    close(inChan)
    done <- 0
}

func chanRunner() {
    doneChan := make(chan int)

    r1 := make(chan gringo.Payload, 256)
    r2 := make(chan gringo.Payload, 256)

    var startTime time.Time = time.Now()

    go chanProducer(r1, doneChan)
    go chanForwarder(r1, r2, doneChan)
    go chanConsumer(r2, doneChan)

    <-doneChan; <-doneChan; <-doneChan

    fmt.Println("chanRunner Nanoseconds passed:", time.Since(startTime))
}

func main() {
    gringoRunner()
    gringoRunner()
    chanRunner()
    chanRunner()
    gringoRunner()
    chanRunner()
}

