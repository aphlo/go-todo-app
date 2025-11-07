package todo

import "encoding/json"

type TodoID struct {
    value int64
}

func NewTodoID(v int64) TodoID { return TodoID{value: v} }
func (id TodoID) Int64() int64 { return id.value }

// json.Marshaler
func (id TodoID) MarshalJSON() ([]byte, error) {
    return json.Marshal(id.value)
}

// json.Unmarshaler
func (id *TodoID) UnmarshalJSON(data []byte) error {
    var v int64
    if err := json.Unmarshal(data, &v); err != nil {
        return err
    }
    id.value = v
    return nil
}
