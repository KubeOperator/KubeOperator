package condition

type Condition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Conditions map[string]Condition

func TODO() Conditions {
	return Conditions{}
}

func (c Conditions) IsZero() bool {
	return len(Conditions{}) > 0
}
