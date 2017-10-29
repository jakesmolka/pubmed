/* Package pubmed provides functions to call PubMed API functions according to
https://www.ncbi.nlm.nih.gov/books/NBK25501/ and
https://www.ncbi.nlm.nih.gov/pmc/tools/id-converter-api/
*/

package pubmed

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

//base urls
var eutilsURL = "http://eutils.ncbi.nlm.nih.gov/entrez/eutils/"
var utilsURL = "https://www.ncbi.nlm.nih.gov/pmc/utils/"

// NewClient creates new API client to be used for all functions
func NewClient() (client *Client, err error) {
	eutils, err := url.ParseRequestURI(eutilsURL)
	if err != nil {
		return
	}
	utils, err := url.ParseRequestURI(utilsURL)
	if err != nil {
		return
	}
	httpClient := http.DefaultClient
	client = &Client{eutils, utils, httpClient}

	return
}

// Esearch submits query as search term and return list of ids etc. of matching articles.
// See: https://www.ncbi.nlm.nih.gov/books/NBK25499/#chapter4.ESearch
func (c *Client) Esearch(query string) (response *EsearchResponse, err error) {
	esearchURL, err := url.ParseRequestURI(c.EutilsURL.String() + "esearch.fcgi?db=pubmed&retmode=json&term=" + query)
	if err != nil {
		return
	}
	body, err := httpGet(esearchURL, c.httpClient)
	if err != nil {
		return
	}

	response = new(EsearchResponse)
	err = json.NewDecoder(body).Decode(&response)
	if err != nil {
		return
	}

	return
}

// Efetch fetches article information for given id in db formatted in retmode.
// Output type depends on db and retmode.
// See: https://www.ncbi.nlm.nih.gov/books/NBK25499/#chapter4.EFetch
// And for possible parameters and combinations: https://www.ncbi.nlm.nih.gov/books/NBK25499/table/chapter4.T._valid_values_of__retmode_and/
func (c *Client) Efetch(id, db, rettype, retmode string) (response interface{}, err error) {
	//efetch := "efetch.fcgi?db=pubmed&retmode=text&rettype=abstract&id="

	// check for wrong input
	if (db != "pubmed") && (db != "pmc") {
		err = errors.New("httpGetJsonEfetch: wrong db")
		return
	}
	if (retmode != "text") && (retmode != "xml") {
		err = errors.New("httpGetJsonEfetch: wrong retmode")
		return
	}

	// build query url and execute
	efetch, err := url.ParseRequestURI(c.EutilsURL.String() + "efetch.fcgi?db=" + db + "&rettype=" + rettype + "&retmode=" + retmode + "&id=" + id)
	if err != nil {
		return nil, err
	}
	body, err := httpGet(efetch, c.httpClient)
	if err != nil {
		return
	}

	// decide processing of response depending on retmode
	if retmode == "text" {
		// retmode=text returns plain text therefore Close() and no decoding
		defer body.Close()
		respByte, e := ioutil.ReadAll(body) //TODO: why is 'err' shadowed here but not below?
		if e != nil {
			return nil, e
		}
		response = string(respByte)
	} else if retmode == "xml" {
		// @xml - decode with struct
		if db == "pubmed" {
			response = new(EfetchPubmedXmlResponse)
		} else if db == "pmc" { // TODO: pmc has default "xml" so retmode="" possible. And rettype is ignored.
			response = new(EfetchPmcXmlResponse)
		}
		err = xml.NewDecoder(body).Decode(&response)
		if err != nil {
			return
		}
	}

	return
}

// IDConvert converts given id between pmid and pmcid. Wraps httpGet with idconv parameters and given id.
func (c *Client) IDConvert(id string) (pmid, pmcid string, err error) {
	// compose URL with needed base URL, formal parameter and given ID
	idconv, err := url.ParseRequestURI(c.UtilsURL.String() + "idconv/v1.0/?format=json&ids=" + id)
	if err != nil {
		return
	}
	body, err := httpGet(idconv, c.httpClient)
	if err != nil {
		return
	}

	response := new(IdconvResponse)
	err = json.NewDecoder(body).Decode(&response)
	if err != nil {
		return "", "", errors.New("IDConvert: failed decoding of response")
	}

	if len(response.Records) == 0 {
		return "", "", errors.New("IDConvert: no records in response")
	}

	pmid = response.Records[0].Pmid
	pmcid = response.Records[0].Pmcid

	return
}

// httpGet is foundation to more specific API functions
func httpGet(requestURL *url.URL, httpClient *http.Client) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close() //TODO: Decode's not working with Close. Is that alright?

	return resp.Body, nil
}
