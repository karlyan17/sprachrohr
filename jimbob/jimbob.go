// jimbob.go
package jimbob

import(
    "encoding/json"
    "reflect"
    "strconv"
    "errors"
    "io/ioutil"
    "sprachrohr/freshlog"
)

type Bucket struct {
    Path string
    Data interface{}
    next_index int
}

func NewBucket(path string, data interface{}) Bucket {
    if validMap(data) {
        bucket := Bucket{
            Path: path,
            Data: data,
        }
        return bucket
    } else {
        freshlog.Fatal.Fatal("Map validity check failed: jimbob bucket data has to be of Type map[int] instead of %v",
        reflect.TypeOf(data),
        )
        return Bucket{}
    }
}

func OpenBucket(path string, elem interface{}) Bucket {
    elem_type := reflect.TypeOf(elem)
    bucket_date_type := reflect.MapOf(reflect.TypeOf(0), elem_type)
    freshlog.Debug.Printf("creating bucket of type %v from elements of type %v", bucket_date_type, elem_type)
    bucket_data := reflect.MakeMap(bucket_date_type)
    doc_files,err := ioutil.ReadDir(path)
    if err != nil {
        freshlog.Fatal.Fatal("Failed to read documents form %v with error %v", path, err)
    }
    for _,file := range(doc_files) {
        freshlog.Debug.Printf("reading file %v", file)
        doc_index,err := strconv.Atoi(file.Name())
        if err != nil {
            freshlog.Error.Printf("could not convert name to index for file %v with error %v", file.Name(), err)
            continue
        }
        contents,err := ioutil.ReadFile(file.Name())
        if err != nil {
            freshlog.Error.Printf("could not read file %v with error %v", file.Name(), err)
            continue
        }
        doc := reflect.New(elem_type)
        err = json.Unmarshal(contents, doc)
        if err != nil {
            freshlog.Error.Printf("could not parse file %v into interface %v with error %v", file.Name(), elem_type, err)
            continue
        }
        bucket_data.SetMapIndex(reflect.ValueOf(doc_index), reflect.ValueOf(doc))

    }
    return NewBucket(path, bucket_data)
}

func (bucket Bucket) AddDoc(doc interface{}, index ...int) (int, error) {
    if len(index) != 1 {
        if _,exists := bucket.Data[index[0]]; !exists {
            bucket.Data[index[0]] = doc
            return index[0], nil
        } else {
            return -1, errors.New(fmt.Sprintf("Index %v already exists in bucket %v", index[0], bucket.Path))
        }
    } else {
        new_index = bucket.next_index
        bucket.data[new_index] = doc
        bucket.updateNextIndex()
        return new_index, nil
    }

}

func (bucket Bucket) updateNextIndex() {
    for i := bucket.next_index;; i++ {
        if _,exists := bucket.data[i]; !exists {
            bucket.next_index = i
        }
    }
}

func validMap(in interface{}) bool {
    if map_kind := reflect.ValueOf(in).Kind().String(); map_kind == "map" {
        if map_keys := reflect.ValueOf(in).MapKeys(); map_keys[0].Kind().String() == "int" {
            freshlog.Debug.Printf("passed map is valid jimbob bucket data")
            return true
        }
    }
    return false
}
