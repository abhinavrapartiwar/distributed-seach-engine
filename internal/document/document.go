package document

type Page struct {
	Title    string   `xml:"title"`
	Revision Revision `xml:"revision"`
}

type Revision struct {
	Text string `xml:"text"`
}

type Document struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
