/* Package pubmed provides functions to call PubMed API functions according to
https://www.ncbi.nlm.nih.gov/books/NBK25501/ and
https://www.ncbi.nlm.nih.gov/pmc/tools/id-converter-api/
*/

package pubmed

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	eutils, err := url.ParseRequestURI(eutilsURL)
	if err != nil {
		t.Error()
	}
	utils, err := url.ParseRequestURI(utilsURL)
	if err != nil {
		t.Error()
	}
	tests := []struct {
		name       string
		wantClient *Client
		wantErr    bool
	}{
		// TODO: Add test cases.
		{"Normal", &Client{eutils, utils, http.DefaultClient}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := NewClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotClient, tt.wantClient) {
				t.Errorf("NewClient() = %v, want %v", gotClient, tt.wantClient)
			}
		})
	}
}

func TestClient_Esearch(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse []string //as in EsearchResult.IdList
		wantErr      bool
	}{
		// TODO: Add test cases.
		{"cancer and retmax 1", args{"cancer&retmax=1&mindate=2017/09/28&maxdate=2017/09/29"}, []string{"28958124"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient()
			if err != nil {
				t.Error("Client.Esearch() error = client creation failed")
			}
			gotResponse, err := c.Esearch(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Esearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse.Result.IdList, tt.wantResponse) {
				t.Errorf("Client.Esearch() = %v, want %v", gotResponse.Result.IdList, tt.wantResponse)
			}
		})
	}
}

func TestClient_Efetch(t *testing.T) {
	type args struct {
		id      string
		db      string
		rettype string
		retmode string
	}
	tests := []struct {
		name         string
		args         args
		wantResponse interface{}
		wantErr      bool
	}{
		// TODO: Add test cases.
		{"pubmed, abstract, text", args{"9997", "pubmed", "abstract", "text"}, `
1. Biochim Biophys Acta. 1976 Sep 28;446(1):179-91.

Magnetic studies of Chromatium flavocytochrome C552. A mechanism for heme-flavin 
interaction.

Strekas TC.

Electron paramagnetic resonance and magnetic susceptibility studies of Chromatium
flavocytochrome C552 and its diheme flavin-free subunit at temperatures below 45 
degrees K are reported. The results show that in the intact protein and the
subunit the two low-spin (S = 1/2) heme irons are distinguishable, giving rise to
separate EPR signals. In the intact protein only, one of the heme irons exists in
two different low spin environments in the pH range 5.5 to 10.5, while the other 
remains in a constant environment. Factors influencing the variable heme iron
environment also influence flavin reactivity, indicating the existence of a
mechanism for heme-flavin interaction.


PMID: 9997  [Indexed for MEDLINE]

`, false},
		{"pubmed, abstract, xml", args{"9997", "pubmed", "abstract", "xml"}, &EfetchPubmedXmlResponse{[]EfetchPmArticle{{"Magnetic studies of Chromatium flavocytochrome C552. A mechanism for heme-flavin interaction.", "Electron paramagnetic resonance and magnetic susceptibility studies of Chromatium flavocytochrome C552 and its diheme flavin-free subunit at temperatures below 45 degrees K are reported. The results show that in the intact protein and the subunit the two low-spin (S = 1/2) heme irons are distinguishable, giving rise to separate EPR signals. In the intact protein only, one of the heme irons exists in two different low spin environments in the pH range 5.5 to 10.5, while the other remains in a constant environment. Factors influencing the variable heme iron environment also influence flavin reactivity, indicating the existence of a mechanism for heme-flavin interaction.", []EfetchPmAuthor{{"Strekas", "T C", "TC"}}, "1976", "Sep", "28"}}}, false},
		//TODO: add case pmc, nil (ignored), xml (which is default). But has huge response!
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient()
			if err != nil {
				t.Error("Client.Efetch() error = client creation failed")
			}
			gotResponse, err := c.Efetch(tt.args.id, tt.args.db, tt.args.rettype, tt.args.retmode)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Efetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("Client.Efetch() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestClient_IDConvert(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name      string
		args      args
		wantPmid  string
		wantPmcid string
		wantErr   bool
	}{
		// TODO: Add test cases.
		{"from example", args{"23193287"}, "23193287", "PMC3531190", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient()
			if err != nil {
				t.Error("Client.IDConvert() error = client creation failed")
			}
			gotPmid, gotPmcid, err := c.IDConvert(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.IDConvert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPmid != tt.wantPmid {
				t.Errorf("Client.IDConvert() gotPmid = %v, want %v", gotPmid, tt.wantPmid)
			}
			if gotPmcid != tt.wantPmcid {
				t.Errorf("Client.IDConvert() gotPmcid = %v, want %v", gotPmcid, tt.wantPmcid)
			}
		})
	}
}

func Test_httpGet(t *testing.T) {
	type args struct {
		requestURL *url.URL
		httpClient *http.Client
	}
	tests := []struct {
		name    string
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := httpGet(tt.args.requestURL, tt.args.httpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("httpGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpGet() = %v, want %v", got, tt.want)
			}
		})
	}
}
