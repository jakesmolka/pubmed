/* Package pubmed provides functions to call PubMed API functions according to
https://www.ncbi.nlm.nih.gov/books/NBK25501/ and
https://www.ncbi.nlm.nih.gov/pmc/tools/id-converter-api/
*/

// TODO: check each type if necessary

package pubmed

import (
	"net/http"
	"net/url"
)

// Client is needed to access API functions. Different URLs for different tool sets.
type Client struct {
	EutilsURL  *url.URL
	UtilsURL   *url.URL
	httpClient *http.Client
}

// data types for rest responses

// -------- esearch --------
type EsearchResult struct {
	Count            string        `json:"count"`
	Retmax           string        `json:"retmax"`
	Retstart         string        `json:"retstart"`
	IdList           []string      `json:"idlist"`
	TranslationSet   []interface{} `json:"translationset"`   //TODO: implement or remove
	TranslationStack []interface{} `json:"translationstack"` //TODO: implement or remove
	QueryTranslation string        `json:"querytranslation"`
}

type EsearchResponse struct {
	Header map[string]string `json:"header"`
	Result EsearchResult     `json:"esearchresult"`
}

// -------- efetch --------
// this respresents response's xml structure using ">" to go through layers
//------ Pubmed
type EfetchPmAuthor struct {
	Lastname string `xml:"LastName"`
	Forename string `xml:"ForeName"`
	Initials string `xml:"Initials"`
}

type EfetchPmArticle struct {
	Title    string           `xml:"MedlineCitation>Article>ArticleTitle"`
	Abstract string           `xml:"MedlineCitation>Article>Abstract>AbstractText"`
	Authors  []EfetchPmAuthor `xml:"MedlineCitation>Article>AuthorList>Author"`
	Year     string           `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Year"`
	Month    string           `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Month"`
	Day      string           `xml:"MedlineCitation>Article>Journal>JournalIssue>PubDate>Day"`
}

type EfetchPubmedXmlResponse struct {
	//XMLName  xml.Name          `xml:"PubmedArticleSet"`
	Articles []EfetchPmArticle `xml:"PubmedArticle"`
}

//------ PMC
type EfetchPmcAuthor struct {
	ContribType string `xml:"contrib-type,attr"` // should be "author" if relevant
	Surename    string `xml:"name>surname"`
	GivenNames  string `xml:"name>given-names"`
}

type EfetchPmcPubdate struct {
	PubType string `xml:"pub-type,attr"` //can be "ppub" and "epub"
	Year    string `xml:"year"`
	Month   string `xml:"month"`
	Day     string `xml:"day"`
}

type EfetchPmcArticle struct {
	Title    string   `xml:"front>article-meta>title-group>article-title"`
	Abstract struct { //want all inner xml
		Data string `xml:",innerxml"`
	} `xml:"front>article-meta>abstract"`
	Authors []EfetchPmcAuthor  `xml:"front>article-meta>contrib-group>contrib"`
	Pubdate []EfetchPmcPubdate `xml:"front>article-meta>pub-date"`
}

type EfetchPmcXmlResponse struct {
	Articles []EfetchPmcArticle `xml:"article"`
}

// -------- idconv --------
type IdconvRecord struct {
	Pmcid  string `json:"pmcid"`
	Pmid   string `json:"pmid"`
	Status string `json:"status"` //allowed? even if field just exists when "error"?
}

type IdconvResponse struct {
	Records []IdconvRecord `json:"records"`
}
