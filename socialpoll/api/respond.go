package main

import(
  ""
)

func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewEncoder(w).Encode(v)
}

func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}