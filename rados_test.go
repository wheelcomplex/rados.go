package rados

import (
    "bytes"
    "fmt"
    "testing"
)

func errorOnError(t *testing.T, e error, message string, parameters ...interface{}) {
    if e != nil {
        t.Errorf("%v : %v", e, fmt.Sprintf(message, parameters...))
    }
}

func fatalOnError(t *testing.T, e error, message string, parameters ...interface{}) {
    if e != nil {
        t.Fatalf("%v : %v", e, fmt.Sprintf(message, parameters...))
    }
}

func Test_RadosNew(t *testing.T) {
    var rados *Rados
    var err error

    rados, err = New("")
    fatalOnError(t, err, "New")

    err = rados.Release()
    fatalOnError(t, err, "Release")

    if rados, err = New("path that does not exist"); err == nil {
        t.Errorf("New should have failed")
        rados.Release()
    }
}

func Test_RadosContext(t *testing.T) {
    var rados *Rados
    var err error

    rados, err = New("")
    fatalOnError(t, err, "New")
    defer rados.Release()

    ctx, err := rados.NewContext("test")
    fatalOnError(t, err, "NewContext")
    ctx.Release()

    if ctx, err = rados.NewContext("pool that does not exist"); err == nil {
        t.Errorf("NewContext should have failed")
        ctx.Release()
    }
}

func Test_RadosObject(t *testing.T) {
    var rados *Rados
    var err error

    rados, err = New("")
    fatalOnError(t, err, "New")
    defer rados.Release()

    ctx, err := rados.NewContext("test")
    fatalOnError(t, err, "NewContext")
    defer ctx.Release()

    name := "test-object"
    name2 := "test-object2"
    data := []byte("test data")

    // Create an object
    _, err = ctx.Create(name)
    errorOnError(t, err, "Create")

    // Make sure it's there
    objInfo, err := ctx.Stat(name)
    fatalOnError(t, err, "Stat")

    if objInfo.Size() != int64(0) {
        t.Errorf("Object size mismatch, was %s, expected %s", objInfo.Size(), 0)
    }

    // Put data in the object
    err = ctx.Put(name, data)
    errorOnError(t, err, "Put")

    // Make sure everything looks right
    objInfo, err = ctx.Stat(name)
    fatalOnError(t, err, "Stat")

    if objInfo.Name() != name {
        t.Errorf("Object name mismatch, was %s, expected %s", objInfo.Name(), name)
    }

    if objInfo.Size() != int64(len(data)) {
        t.Errorf("Object size mismatch, was %d, expected %d", objInfo.Size(), len(data))
    }

    // Get the data back
    data2, err := ctx.Get(name)
    fatalOnError(t, err, "Get")

    // It better be the same
    if !bytes.Equal(data, data2) {
        t.Errorf("Object data mismatch, was %s, expected %s", data2, data)
    }

    // Open an existing object
    obj, err := ctx.Open(name)
    fatalOnError(t, err, "Open")

    // Make sure everything looks right
    if obj.Name() != name {
        t.Errorf("Object name mismatch, was %s, expected %s", obj.Name(), name)
    }

    if obj.Size() != int64(len(data)) {
        t.Errorf("Object size mismatch, was %d, expected %d", obj.Size(), len(data))
    }

    // Open a new object
    obj, err = ctx.Open(name2)
    fatalOnError(t, err, "Open")

    // Make sure it's there
    objInfo, err = ctx.Stat(name2)
    fatalOnError(t, err, "Stat")

    if objInfo.Size() != int64(0) {
        t.Errorf("Object size mismatch, was %d, expected %d", objInfo.Size(), 0)
    }

    // Remove the objects
    err = ctx.Remove(name)
    errorOnError(t, err, "Remove")
    err = ctx.Remove(name2)
    errorOnError(t, err, "Remove")

    // They should be gone
    objInfo, err = ctx.Stat(name)
    if err == nil {
        t.Errorf("Object %s should have been deleted be status returned success", name)
    }

    objInfo, err = ctx.Stat(name2)
    if err == nil {
        t.Errorf("Object %s should have been deleted be status returned success", name2)
    }
}