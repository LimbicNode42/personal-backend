package helpers

// import (

// )

func RemoveByValue(slice []*string, value *string) []*string {
    for i, v := range slice {
        if v != nil && value != nil && *v == *value {
            return append(slice[:i], slice[i+1:]...)
        }
    }
    return slice
}
