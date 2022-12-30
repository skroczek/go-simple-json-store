package fs

import (
	"fmt"
	"github.com/skroczek/acme-restful/pkg/backend"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// emptyTree is an empty tree
var emptyTree = make(map[string]interface{})

/*
singleLevelTree is a tree with a single level
/
└─ foo.json
*/
var singleLevelTree = map[string]interface{}{
	"foo.json": &Blob{
		Content: []byte("{name: \"foo\"}"),
		ModTime: time.Now(),
	},
}

/*
multilevelTree is a tree with multiple levels

/
├─ foo/
│ ├─ bar.json
│ └─ baz.json
*/
var multilevelTree = map[string]interface{}{
	"foo": map[string]interface{}{
		"bar.json": &Blob{Content: []byte("{name: \"baz\"}"), ModTime: time.Now()},
		"baz.json": &Blob{Content: []byte("{name: \"qux\"}"), ModTime: time.Now()},
	},
}

func TestMemory_Exists(t *testing.T) {
	type fields struct {
		tree map[string]interface{}
	}
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// test invalid path no .json
		{
			name: "test invalid path /.json",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: "/.json",
			},
			want:    false,
			wantErr: true,
		},
		// test invalid path no .json
		{
			name: "test invalid path no .json",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: "invalid",
			},
			want:    false,
			wantErr: true,
		},
		// test path to short
		{
			name: "test path to short",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: ".json",
			},
			want:    false,
			wantErr: true,
		},
		// test path not found
		{
			name: "test path not found",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: "foo.json",
			},
			want:    false,
			wantErr: false,
		},
		// test path found single level
		{
			name: "test path found single level",
			fields: fields{
				tree: singleLevelTree,
			},
			args: args{
				path: "foo.json",
			},
			want:    true,
			wantErr: false,
		},
		// test path not found single level
		{
			name: "test path not found single level",
			fields: fields{
				tree: singleLevelTree,
			},
			args: args{
				path: "foo1.json",
			},
			want:    false,
			wantErr: false,
		},
		// test path found multiple levels
		{
			name: "test path found multiple levels",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "foo/bar.json",
			},
			want:    true,
			wantErr: false,
		},
		// test path not found multiple levels
		{
			name: "test path not found multiple levels",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "foo/bar2.json",
			},
			want:    false,
			wantErr: false,
		},
		// test folder not found multiple levels
		{
			name: "test folder not found multiple levels",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "foo/bar/bar.json",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				tree: tt.fields.tree,
			}
			got, err := m.Exists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_Get(t *testing.T) {
	type fields struct {
		tree map[string]interface{}
	}
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// test invalid path no .json
		{
			name: "test invalid path no .json",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		// test path to short
		{
			name: "test path to short",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: ".json",
			},
			want:    nil,
			wantErr: true,
		},
		// test path not found
		{
			name: "test path not found",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: "foo.json",
			},
			want:    nil,
			wantErr: true,
		},
		// test path found single level
		{
			name: "test path found single level",
			fields: fields{
				tree: singleLevelTree,
			},
			args: args{
				path: "foo.json",
			},
			want:    []byte("{name: \"foo\"}"),
			wantErr: false,
		},
		// test path not found single level
		{
			name: "test path not found single level",
			fields: fields{
				tree: singleLevelTree,
			},
			args: args{
				path: "foo1.json",
			},
			want:    nil,
			wantErr: true,
		},
		// test path found multiple levels
		{
			name: "test path found multiple levels",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "foo/bar.json",
			},
			want:    []byte("{name: \"baz\"}"),
			wantErr: false,
		},
		// test path not found multiple levels
		{
			name: "test path not found multiple levels",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "foo/bar2.json",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				tree: tt.fields.tree,
			}
			got, err := m.Get(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_List(t *testing.T) {
	type fields struct {
		tree map[string]interface{}
	}
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// test invalid path no .json
		{
			name: "test invalid path no .json",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
		// test path to short
		{
			name: "test path to short",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: ".json",
			},
			want:    nil,
			wantErr: true,
		},
		// test path not found
		{
			name: "test path not found",
			fields: fields{
				tree: emptyTree,
			},
			args: args{
				path: "foo.json",
			},
			want:    nil,
			wantErr: true,
		},
		// test path found single level
		{
			name: "test path found single level",
			fields: fields{
				tree: singleLevelTree,
			},
			args: args{
				path: "",
			},
			want:    []string{"foo.json"},
			wantErr: false,
		},
		// test path not found single level
		{
			name: "test path not found single level",
			fields: fields{
				tree: singleLevelTree,
			},
			args: args{
				path: "foo1",
			},
			want:    nil,
			wantErr: true,
		},
		// test path found multiple levels
		{
			name: "test path found multiple levels",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "foo",
			},
			want:    []string{"bar.json", "baz.json"},
			wantErr: false,
		},
		// test path found multiple levels
		{
			name: "test path found multiple levels tailing slash",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "/foo/",
			},
			want:    []string{"bar.json", "baz.json"},
			wantErr: false,
		},
		// test path not found multiple levels
		{
			name: "test path not found multiple levels",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "foo/bar2",
			},
			want:    nil,
			wantErr: true,
		},
		// test dir not found multiple levels
		{
			name: "test dir not found multiple levels",
			fields: fields{
				tree: multilevelTree,
			},
			args: args{
				path: "foo/bar/bar2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				tree: tt.fields.tree,
			}
			got, err := m.List(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Slice(got, func(i, j int) bool {
				return got[i] < got[j]
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i] < tt.want[j]
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func NewTree(path string, content *Blob) map[string]interface{} {
	tree := make(map[string]interface{})
	if path == "" {
		return tree
	}
	parts := strings.Split(path, "/")
	runnable := tree
	for i := 0; i < len(parts)-1; i++ {
		runnable[parts[i]] = make(map[string]interface{})
		runnable = runnable[parts[i]].(map[string]interface{})
	}
	runnable[parts[len(parts)-1]] = content
	fmt.Printf("tree: %v\n", tree)
	return tree
}

func TestMemory_Write(t *testing.T) {
	type fields struct {
		tree map[string]interface{}
	}
	type args struct {
		path string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		// test invalid path no .json
		{
			name: "test invalid path no .json",
			args: args{
				path: "invalid",
				data: []byte("test"),
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			wantErr: true,
		},
		// test path to short
		{
			name: "test path to short",
			args: args{
				path: ".json",
				data: []byte("test"),
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			wantErr: true,
		},
		// test path found single level
		{
			name: "test path found single level",
			args: args{
				path: "foo.json",
				data: []byte("test"),
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			wantErr: false,
		},
		// test path found multiple levels
		{
			name: "test path found multiple levels",
			args: args{
				path: "foo/bar/baz.json",
				data: []byte("test"),
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			wantErr: false,
		},
		// test path found multiple levels tailing slash
		{
			name: "test path multiple levels tailing slash",
			args: args{
				path: "foo/bar/baz.json/",
				data: []byte("test"),
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			wantErr: true, // because of the tailing slash it is not a valid path
		},
		// test path found multiple levels tailing slash
		{
			name: "test path already exist multiple levels",
			args: args{
				path: "foo/bar/baz.json",
				data: []byte("test"),
			},
			fields: fields{
				tree: NewTree("foo/bar/baz.json", &Blob{
					Content: []byte("test1"),
					ModTime: time.Now(),
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				tree: tt.fields.tree,
			}
			if err := m.Write(tt.args.path, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Delete(t *testing.T) {
	type fields struct {
		tree map[string]interface{}
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		fields  fields
		wantErr bool
	}{
		// test path not found single level
		{
			name: "test path not found single level",
			args: args{
				path: "foo.json",
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			wantErr: true,
		},
		// test path not found multiple levels
		{
			name: "test path not found multiple levels",
			args: args{
				path: "foo/bar/baz.json",
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			wantErr: true,
		},
		// test path found single level
		{
			name: "test path found single level",
			args: args{
				path: "foo.json",
			},
			fields: fields{
				tree: NewTree("foo.json", &Blob{
					Content: []byte("test1"),
					ModTime: time.Now(),
				}),
			},
			wantErr: false,
		},
		// test path found multiple levels
		{
			name: "test path found multiple levels",
			args: args{
				path: "foo/bar/baz.json",
			},
			fields: fields{
				tree: NewTree("foo/bar/baz.json", &Blob{
					Content: []byte("test1"),
					ModTime: time.Now(),
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				tree: tt.fields.tree,
			}
			if err := m.Delete(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_GetLastModified(t *testing.T) {
	type fields struct {
		tree map[string]interface{}
	}
	type args struct {
		path string
	}
	now := time.Now()
	tests := []struct {
		name    string
		args    args
		fields  fields
		want    time.Time
		wantErr bool
	}{
		// test path not found single level
		{
			name: "test path not found single level",
			args: args{
				path: "foo.json",
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			want:    time.Time{},
			wantErr: true,
		},
		// test path not found multiple levels
		{
			name: "test path not found multiple levels",
			args: args{
				path: "foo/bar/baz.json",
			},
			fields: fields{
				tree: make(map[string]interface{}),
			},
			want:    time.Time{},
			wantErr: true,
		},
		// test path found single level
		{
			name: "test path found single level",
			args: args{
				path: "foo.json",
			},
			fields: fields{
				tree: NewTree("foo.json", &Blob{
					Content: []byte("test1"),
					ModTime: now,
				}),
			},
			want:    now,
			wantErr: false,
		},
		// test path found multiple levels
		{
			name: "test path found multiple levels",
			args: args{
				path: "foo/bar/baz.json",
			},
			fields: fields{
				tree: NewTree("foo/bar/baz.json", &Blob{
					Content: []byte("test1"),
					ModTime: now,
				}),
			},
			want:    now,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				tree: tt.fields.tree,
			}
			got, err := m.GetLastModified(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastModified() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastModified() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMemory(t *testing.T) {
	tests := []struct {
		name string
		want backend.Backend
	}{
		{
			name: "test new memory",
			want: &Memory{
				tree: make(map[string]interface{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemory() = %v, want %v", got, tt.want)
			}
		})
	}
}
