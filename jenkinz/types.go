package jenkinz

type Jobs struct {
	Class string `json:"_class"`
	Jobs  []struct {
		Name string `json:"name"`
	} `json:"jobs"`
}
type Builds struct {
	Builds []struct {
		ID string `json:"id"`
	} `json:"builds"`
}

type Build struct {
	Job string
	Id  string
}
