package main

func main() {
	//rules.RunRules()
	//files.ReadExcelData("C:\\Users\\Srav\\Documents\\code\\irules\\R.xlsx")

}

/*
type Person struct {
	Name    string
	Age     int
	IsAlive bool
	Address Address
}

type Address struct {
	Street string
	City   string
	State  string
}

func printStructValues(v interface{}) {
	val := reflect.ValueOf(v)
	printValues(val, 0)
}

func printValues(val reflect.Value, level int) {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fmt.Printf("%s%s: ", indent(level), val.Type().Field(i).Name)
			if field.Kind() == reflect.Struct {
				fmt.Println()
				printValues(field, level+1)
			} else {
				fmt.Printf("%v\n", field.Interface())
			}
		}
	}
}

func indent(level int) string {
	return string(make([]rune, level*2))
}

func mainxx() {
	address := Address{
		Street: "123 Elm St",
		City:   "Somewhere",
		State:  "CA",
	}

	person := Person{
		Name:    "John Doe",
		Age:     30,
		IsAlive: true,
		Address: address,
	}

	printStructValues(person)
}
*/
