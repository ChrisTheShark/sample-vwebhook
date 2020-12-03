package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	. "github.com/ChrisTheShark/simple-vwebhook"
	"github.com/stretchr/testify/assert"
)

func TestHappyPath(t *testing.T) {
	ns := Namespace{
		Metadata: Metadata{
			Name: "test-ns",
			Labels: map[string]string{
				"team": "avengers",
			},
		},
	}

	bs, _ := json.Marshal(ns)
	admReview := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: bs,
			},
		},
	}

	bs2, _ := json.Marshal(admReview)
	r := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(bs2))
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()

	bs3, _ := ioutil.ReadAll(resp.Body)
	rcvdAdmReview := v1beta1.AdmissionReview{}
	json.Unmarshal(bs3, &rcvdAdmReview)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, true, rcvdAdmReview.Response.Allowed)
}

func TestNoBody(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/validate", nil)
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestInvalidBody(t *testing.T) {
	fake := struct {
		Name string
	}{
		"jon",
	}

	bs, _ := json.Marshal(fake)
	r := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(bs))
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestMissingRequiredLabel(t *testing.T) {
	ns := Namespace{
		Metadata: Metadata{
			Name:   "test-ns",
			Labels: map[string]string{},
		},
	}

	bs, _ := json.Marshal(ns)
	admReview := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: bs,
			},
		},
	}

	bs2, _ := json.Marshal(admReview)
	r := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(bs2))
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()

	bs3, _ := ioutil.ReadAll(resp.Body)
	rcvdAdmReview := v1beta1.AdmissionReview{}
	json.Unmarshal(bs3, &rcvdAdmReview)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, false, rcvdAdmReview.Response.Allowed)
	assert.Equal(t, InvalidMessage, rcvdAdmReview.Response.Result.Message)
}

func TestMissingRequiredLabel2(t *testing.T) {
	ns := Namespace{
		Metadata: Metadata{
			Name: "test-ns",
		},
	}

	bs, _ := json.Marshal(ns)
	admReview := v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: bs,
			},
		},
	}

	bs2, _ := json.Marshal(admReview)
	r := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(bs2))
	w := httptest.NewRecorder()

	Validate(w, r)
	resp := w.Result()

	bs3, _ := ioutil.ReadAll(resp.Body)
	rcvdAdmReview := v1beta1.AdmissionReview{}
	json.Unmarshal(bs3, &rcvdAdmReview)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, false, rcvdAdmReview.Response.Allowed)
	assert.Equal(t, InvalidMessage, rcvdAdmReview.Response.Result.Message)
}
