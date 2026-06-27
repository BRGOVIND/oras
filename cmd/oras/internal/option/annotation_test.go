/*
Copyright The ORAS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package option

import (
	"errors"
	"reflect"
	"testing"
)

func TestAnnotation_Parse_manifestDefault(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{"key=value"},
	}
	if err := opts.Parse(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]map[string]string{
		AnnotationManifest: {"key": "value"},
	}
	if !reflect.DeepEqual(opts.Annotations, want) {
		t.Errorf("got %v, want %v", opts.Annotations, want)
	}
}

func TestAnnotation_Parse_explicitManifest(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{"$manifest:key=value"},
	}
	if err := opts.Parse(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]map[string]string{
		AnnotationManifest: {"key": "value"},
	}
	if !reflect.DeepEqual(opts.Annotations, want) {
		t.Errorf("got %v, want %v", opts.Annotations, want)
	}
}

func TestAnnotation_Parse_config(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{"$config:key=value"},
	}
	if err := opts.Parse(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]map[string]string{
		AnnotationConfig: {"key": "value"},
	}
	if !reflect.DeepEqual(opts.Annotations, want) {
		t.Errorf("got %v, want %v", opts.Annotations, want)
	}
}

func TestAnnotation_Parse_fileLayer(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{"hi.txt:key=value"},
	}
	if err := opts.Parse(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]map[string]string{
		"hi.txt": {"key": "value"},
	}
	if !reflect.DeepEqual(opts.Annotations, want) {
		t.Errorf("got %v, want %v", opts.Annotations, want)
	}
}

func TestAnnotation_Parse_mixed(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{
			"mkey=mval",
			"$config:ckey=cval",
			"hi.txt:lkey=lval",
		},
	}
	if err := opts.Parse(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]map[string]string{
		AnnotationManifest: {"mkey": "mval"},
		AnnotationConfig:   {"ckey": "cval"},
		"hi.txt":           {"lkey": "lval"},
	}
	if !reflect.DeepEqual(opts.Annotations, want) {
		t.Errorf("got %v, want %v", opts.Annotations, want)
	}
}

func TestAnnotation_Parse_emptyInput(t *testing.T) {
	opts := Annotation{}
	if err := opts.Parse(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(opts.Annotations) != 0 {
		t.Errorf("expected empty annotations map, got %v", opts.Annotations)
	}
}

func TestAnnotation_Parse_missingEquals(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{"badformat"},
	}
	if err := opts.Parse(nil); !errors.Is(err, errAnnotationFormat) {
		t.Errorf("expected errAnnotationFormat, got %v", err)
	}
}

func TestAnnotation_Parse_duplicateKeyManifest(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{"key=val1", "key=val2"},
	}
	if err := opts.Parse(nil); !errors.Is(err, errAnnotationDuplication) {
		t.Errorf("expected errAnnotationDuplication, got %v", err)
	}
}

func TestAnnotation_Parse_duplicateKeyConfig(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{"$config:key=val1", "$config:key=val2"},
	}
	if err := opts.Parse(nil); !errors.Is(err, errAnnotationDuplication) {
		t.Errorf("expected errAnnotationDuplication, got %v", err)
	}
}

func TestAnnotation_Parse_duplicateKeyFile(t *testing.T) {
	opts := Annotation{
		ManifestAnnotations: []string{"hi.txt:key=val1", "hi.txt:key=val2"},
	}
	if err := opts.Parse(nil); !errors.Is(err, errAnnotationDuplication) {
		t.Errorf("expected errAnnotationDuplication, got %v", err)
	}
}
