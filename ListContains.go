package main 

func ListContains(stringList[]string, stringCheck string)bool{
	for _, i := range stringList{ 
		if i == stringCheck{
			return true
		}
	}
	return false
}