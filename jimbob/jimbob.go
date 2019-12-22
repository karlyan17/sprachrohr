// jimbob.go
package jimbob

import(
    "encoding/json"
    "strconv"
    "errors"
    "fmt"
    "io/ioutil"
)

type Bucket struct {
    Path string
    Data map[int]interface{}
    next_index int
}

func NewBucket(path string) (Bucket,error) {
    err := validDir(path)
    if err != nil {
        return Bucket{},err
    } else {
        bucket := Bucket{
            Path: path,
            Data: make(map[int]interface{}),
        }
        bucket.updateNextIndex()
        return bucket,nil
    }
}

func OpenBucket(path string) (Bucket,error) {
    err := validDir(path)
    if err != nil {
        return Bucket{},err
    }
    doc_files,err := ioutil.ReadDir(path)
    if err != nil {
        return Bucket{},err
    }

    var return_err error
    bucket := Bucket{
        Path: path,
        Data: make(map[int]interface{}),
    }

    for _,file := range(doc_files) {
        var doc interface{}
        doc_index,err := strconv.Atoi(file.Name())
        if err != nil {
            return_err = errors.New(fmt.Sprintf("%v\ncould not convert name to index for file %v with error %v", return_err, file.Name(), err))
            continue
        }
        path_name := path + "/" + file.Name()
        contents,err := ioutil.ReadFile(path_name)
        if err != nil {
            return_err = errors.New(fmt.Sprintf("%v\ncould not read file %v with error %v", return_err, path_name, err))
            continue
        }

        err = json.Unmarshal(contents, &doc)
        if err != nil {
            return_err = errors.New(fmt.Sprintf("%v\ncould not parse file %v  with error %v", return_err, path_name, err))
            continue
        }
        bucket.Data[doc_index] = doc

    }
    bucket.updateNextIndex()
    return bucket, return_err
}

func Commit(path string,data map[int]interface{}) error {
    var ret_err error
    var path_name string
    for name, doc := range(data) {
        doc_json,err := json.Marshal(doc)
        if err != nil {
            ret_err  = errors.New(fmt.Sprintf("%v\nunable to enjode %v:%v into json skipping",ret_err,doc,name))
            continue
        }
        path_name = path + "/" + strconv.Itoa(name)
        err = ioutil.WriteFile(path_name, doc_json, 0644)
        if err != nil {
            ret_err = errors.New(fmt.Sprintf("%v\nunable to write file ",ret_err,path_name))
            continue
        }
    }
    return ret_err
}

func (bucket *Bucket) Post(doc interface{}) (int, error) {
    new_index := bucket.next_index
    if _,exists := bucket.Data[new_index]; exists {
        return -1, errors.New(fmt.Sprintf("entry with index %v already exists:[%v], this means the database is in an undefined state, consider restarting application", new_index, bucket.Data[new_index]))
    } else {
        bucket.Data[new_index] = doc
        bucket.updateNextIndex()
        go Commit(bucket.Path, bucket.Data)
        return new_index, nil
    }
}

func emptyDir(path string) error {
    dir_cont,err := ioutil.ReadDir(path)
    if err != nil {
        return err
    } else if len(dir_cont) > 0 {
        return errors.New(fmt.Sprintf("%v is not empty, may already contain a jimbob bucket. Not importing", path))
    } else {
        return nil
    }
}

func validDir(path string) error {
    dir_cont,err := ioutil.ReadDir(path)
    if err != nil {
        return err
    }
    for _,file := range(dir_cont) {
        if file.IsDir() {
            return errors.New(fmt.Sprintf("%v contains the directory %v. jimbob filestructure does not use directories. Not importing", path, file.Name()))
        }
    }
    return nil
}

func (bucket *Bucket) updateNextIndex() {
    for i := bucket.next_index;; i++ {
        if _,exists := bucket.Data[i]; !exists {
            bucket.next_index = i
            break
        }
    }
}
