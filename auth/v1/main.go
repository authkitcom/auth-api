package main

import (
	"fmt"
	"gitlab.authkit.com/authkit/core/gen"
	"os"
	"text/template"
)

//go:generate go build
//go:generate ./v1

type ServiceOperationTemplate struct {
	ProtoTemplate     string
	ProtoTypeTemplate string
}

type MainInput struct {
	MainServiceName string
	ResourceName    string
	Services        []*ServiceInput
}

type ServiceOperation struct {
	Input  *ServiceInput
	Name   string
	Params map[string]interface{}
}

type ServiceInput struct {
	Name         string
	Type         string
	Resource     string
	Operations   []*ServiceOperation
	TenantScoped bool
	RealmScoped  bool
}

var operationTemplates = map[string]*ServiceOperationTemplate{
	"Associate": {
		ProtoTemplate: `
  rpc Associate{{ .Params.Field  }}sTo{{ .Input.Name }} (Associate{{ .Params.Field  }}sTo{{ .Input.Name }}Request) returns (Associate{{ .Params.Field  }}sTo{{ .Input.Name }}Response) {
    option (google.api.http) = {
      put: "/v1/{{ if .Input.TenantScoped}}{tenant_scope}/{{end}}{{ if .Input.RealmScoped}}{realm_scope}/{{end}}{{ .Input.Resource }}s/{{ "{" }}{{ .Input.Resource }}{{ "}" }}/{{ .Params.Resource }}"
      body: "{{ .Params.ParamName }}"
    };

    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
	  security: {
	    security_requirement: {
		  key: "OAuth2";
		    value: {
		      scope: "authkit.com/auth:all";
			}
        }
	  }
	};
  }
`,
		ProtoTypeTemplate: `
message Associate{{ .Params.Field  }}sTo{{ .Input.Name }}Request {
  {{ .Input.Name }}{{ .Params.Field }}Association {{ .Params.ParamName }} = 1;
  string {{ .Input.Resource }} = 2;
{{- if .Input.TenantScoped }}
  string tenant_scope = 3;
{{- end }}
{{- if .Input.RealmScoped }}
  string realm_scope = 4;
{{- end }}
}

message Associate{{ .Params.Field  }}sTo{{ .Input.Name }}Response {
}

message {{ .Input.Name }}{{ .Params.Field }}Association {
  repeated string set = 1;
  repeated string remove = 2;
}
`,
	},
	"ListAssociation": {
		ProtoTemplate: `
  rpc List{{ .Params.Field  }}sBy{{ .Input.Name }} (List{{ .Params.Field  }}sBy{{ .Input.Name }}Request) returns (List{{ .Params.Field  }}sBy{{ .Input.Name }}Response) {
    option (google.api.http) = {
	  get: "/v1/{{ if .Input.TenantScoped}}{tenant_scope}/{{end}}{{ if .Input.RealmScoped}}{realm_scope}/{{end}}{{ .Input.Resource }}s/{{ "{" }}{{ .Input.Resource }}{{ "}" }}/{{ .Params.Resource }}"
    };

    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
	  security: {
	    security_requirement: {
		  key: "OAuth2";
		    value: {
		      scope: "authkit.com/auth:all";
			}
        }
	  }
	};
  }
		`,
		ProtoTypeTemplate: `
message List{{ .Params.Field  }}sBy{{ .Input.Name }}Request {
  string sorting = 1;
  auth.v1.PageParams paging = 2;
  string {{ .Input.Resource }} = 3;
{{- if .Input.TenantScoped }}
  string tenant_scope = 4;
{{- end }}
{{- if .Input.RealmScoped }}
  string realm_scope = 5;
{{- end }}
}

message List{{ .Params.Field  }}sBy{{ .Input.Name }}Response {
  auth.v1.PageInfo page_info = 1;
  repeated {{ .Params.Type }} list = 2;
}
		`,
	},
	"Find": {
		ProtoTemplate: `
  rpc Get{{ .Input.Name }}(Get{{ .Input.Name }}Request) returns (Get{{ .Input.Name }}Response) {
    option (google.api.http) = { get: "/v1/{{ if .Input.TenantScoped}}{tenant_scope}/{{end}}{{ if .Input.RealmScoped}}{realm_scope}/{{end}}{{ .Input.Resource }}s/{{ "{" }}{{ .Input.Resource }}{{ "}" }}" };
  }
`,
		ProtoTypeTemplate: `
message Get{{ .Input.Name }}Request {
  string {{ .Input.Resource }} = 1;
{{- if .Input.TenantScoped }}
  string tenant_scope = 2;
{{- end }}
{{- if .Input.RealmScoped }}
  string realm_scope = 3;
{{- end }}
}

message Get{{ .Input.Name }}Response {
  {{ .Input.Type }} {{ .Input.Resource }} = 1;
}
`,
	},
	"CreateUpdateDelete": {
		ProtoTemplate: `
  rpc Create{{ .Input.Name }}(Create{{ .Input.Name }}Request) returns (Create{{ .Input.Name }}Response) {
    option (google.api.http) = {
      post: "/v1/{{ if .Input.TenantScoped}}{tenant_scope}/{{end}}{{ if .Input.RealmScoped}}{realm_scope}/{{end}}{{ .Input.Resource }}s"
      body: "{{ .Input.Resource }}"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
	  security: {
	    security_requirement: {
		  key: "OAuth2";
		    value: {
		      scope: "authkit.com/auth:all";
			}
        }
	  }
	};
  }

  rpc Update{{ .Input.Name }}(Update{{ .Input.Name }}Request) returns (Update{{ .Input.Name }}Response) {
    option (google.api.http) = {
      put: "/v1/{{ if .Input.TenantScoped}}{tenant_scope}/{{end}}{{ if .Input.RealmScoped}}{realm_scope}/{{end}}{{ .Input.Resource }}s/{ 
{{- .Input.Resource }}.id}"
      body: "{{ .Input.Resource }}"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
	  security: {
	    security_requirement: {
		  key: "OAuth2";
		    value: {
		      scope: "authkit.com/auth:all";
			}
        }
	  }
	};
  }

  rpc Delete{{ .Input.Name }}(Delete{{ .Input.Name }}Request) returns (Delete{{ .Input.Name }}Response) {
    option (google.api.http) = { delete: "/v1/{{ if .Input.TenantScoped}}{tenant_scope}/{{end}}{{ if .Input.RealmScoped}}{realm_scope}/{{end}}{{ .Input.Resource }}s/{id}" };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
	  security: {
	    security_requirement: {
		  key: "OAuth2";
		    value: {
		      scope: "authkit.com/auth:all";
			}
        }
	  }
	};
  }
`,
		ProtoTypeTemplate: `
message Create{{ .Input.Name }}Request {
  {{ .Input.Type }} {{ .Input.Resource }} = 1;
{{- if .Input.TenantScoped }}
  string tenant_scope = 2;
{{- end }}
{{- if .Input.RealmScoped }}
  string realm_scope = 3;
{{- end }}
}

message Create{{ .Input.Name }}Response {
  {{ .Input.Type }} {{ .Input.Resource }} = 1;
}

message Update{{ .Input.Name }}Request {
  {{ .Input.Type }} {{ .Input.Resource }} = 1;
{{- if .Input.TenantScoped }}
  string tenant_scope = 3;
{{- end }}
{{- if .Input.RealmScoped }}
  string realm_scope = 4;
{{- end }}
}

message Update{{ .Input.Name }}Response {
  {{ .Input.Type }} {{ .Input.Resource }} = 1;
}

message Delete{{ .Input.Name }}Request {
  string id = 1;
{{- if .Input.TenantScoped }}
  string tenant_scope = 2;
{{- end }}
{{- if .Input.RealmScoped }}
  string realm_scope = 3;
{{- end }}
}

message Delete{{ .Input.Name }}Response {
}

`,
	},
	"UpdateConfig": {
		ProtoTemplate: `
  rpc Update{{ .Input.Name }}Config(Update{{ .Input.Name }}ConfigRequest) returns (Update{{ .Input.Name }}ConfigResponse) {
    option (google.api.http) = {
      put: "/v1/{{ if .Input.TenantScoped}}{tenant_scope}/{{end}}{{ if .Input.RealmScoped}}{realm_scope}/{{end}}{{ .Input.Resource }}s/{id}/config"
      body: "config"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
	  security: {
	    security_requirement: {
		  key: "OAuth2";
		    value: {
		      scope: "authkit.com/auth:all";
			}
        }
	  }
	};
  }
`,
		ProtoTypeTemplate: `
message Update{{ .Input.Name }}ConfigRequest {
  Update{{ .Input.Name }}Config config = 1;
  string id = 2;
{{- if .Input.TenantScoped }}
  string tenant_scope = 3;
{{- end }}
{{- if .Input.RealmScoped }}
  string realm_scope = 4;
{{- end }}
}

message Update{{ .Input.Name }}Config {
  google.protobuf.Struct set = 1;
  repeated string remove = 2;
}

message Update{{ .Input.Name }}ConfigResponse {
}
`,
	},
	"List": {
		ProtoTemplate: `
  rpc List{{ .Input.Name }}s(List{{ .Input.Name }}sRequest) returns (List{{ .Input.Name }}sResponse) {
    option (google.api.http) = { get: "/v1/{{ if .Input.TenantScoped}}{tenant_scope}/{{end}}{{ if .Input.RealmScoped}}{realm_scope}/{{end}}{{ .Input.Resource }}s" };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
	  security: {
	    security_requirement: {
		  key: "OAuth2";
		    value: {
		      scope: "authkit.com/auth:all";
			}
        }
	  }
	};
  }
`,
		ProtoTypeTemplate: `
message List{{ .Input.Name }}sRequest {
  string sorting = 1;
  auth.v1.PageParams paging = 2;
{{- if .Input.TenantScoped }}
  string tenant_scope = 3;
{{- end }}
{{- if .Input.RealmScoped }}
  string realm_scope = 4;
{{- end }}
}

message List{{ .Input.Name }}sResponse {
  auth.v1.PageInfo page_info = 1;
  repeated {{ .Input.Type }} list = 2;
}
`,
	},
}

var protoTemplate = `// GENERATED BY go:generate. DO NOT EDIT.

syntax = "proto3";
package auth.v1;

import "google/api/annotations.proto";
import "auth/v1/auth.proto";
import "google/protobuf/struct.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
	info: {
		title: "AuthKit API";
		version: "0.1";
		contact: {
			name: "AuthKit Support";
			url: "https://support.authkit.com";
			email: "support@authkit.com";
		};
	};
	security_definitions: {
		security: {
			key: "OAuth2";
			value: {
				type: TYPE_OAUTH2;
				flow: FLOW_ACCESS_CODE;
				authorization_url: "{{ "{{" }} .IDPIssuer {{ "}}" }}/authorize";
				token_url: "{{ "{{" }} .IDPIssuer {{ "}}" }}/oauth/token";
				scopes: {
					scope: {
						key: "authkit.com/auth:all";
						value: "All API access";
					}
				}
			}
		}
	}
	security: {
		security_requirement: {
			key: "OAuth2";
			value: {
				scope: "authkit.com/auth:all";
			}
		}
	}
};

service {{ .MainServiceName }} {

{{- range .Services }}
{{- range .Operations }}
  {{ CallTemplate (print .Name "Proto") . }}
{{- end }}
{{- end }}
}

{{- range .Services }}
{{- range .Operations }}
  {{ CallTemplate (print .Name "ProtoType") . }}
{{- end }}
{{- end }}

`

func MakeServiceInput(input *ServiceInput, operations []*ServiceOperation) *ServiceInput {

	for _, v := range operations {
		v.Input = input
	}

	input.Operations = operations

	return input
}

func main() {

	input := []*MainInput{
		{
			MainServiceName: "AuthKitAuthService",
			ResourceName:    "auth",
			Services: []*ServiceInput{
				MakeServiceInput(&ServiceInput{
					Name:     "Tenant",
					Type:     "auth.v1.Tenant",
					Resource: "tenant",
				}, []*ServiceOperation{
					{Name: "List"},
				}),
				MakeServiceInput(&ServiceInput{
					Name:         "User",
					Type:         "auth.v1.User",
					Resource:     "user",
					TenantScoped: true,
				}, []*ServiceOperation{
					{Name: "Find"},
					{Name: "List"},
					{Name: "CreateUpdateDelete"},
				}),
				MakeServiceInput(&ServiceInput{
					Name:         "UserDatabase",
					Type:         "auth.v1.UserDatabase",
					Resource:     "user_database",
					TenantScoped: true,
				}, []*ServiceOperation{
					{Name: "Find"},
					{Name: "List"},
					{Name: "CreateUpdateDelete"},
					{Name: "UpdateConfig"},
				}),
				MakeServiceInput(&ServiceInput{
					Name:         "Provider",
					Type:         "auth.v1.Provider",
					Resource:     "provider",
					TenantScoped: true,
				}, []*ServiceOperation{
					{Name: "Find"},
					{Name: "List"},
					{Name: "CreateUpdateDelete"},
					{Name: "UpdateConfig"},
				}),
				MakeServiceInput(&ServiceInput{
					Name:         "Realm",
					Type:         "auth.v1.Realm",
					Resource:     "realm",
					TenantScoped: true,
				}, []*ServiceOperation{
					{Name: "Find"},
					{Name: "List"},
					{Name: "CreateUpdateDelete"},
				}),
				MakeServiceInput(&ServiceInput{
					Name:         "Client",
					Type:         "auth.v1.Client",
					Resource:     "client",
					TenantScoped: true,
					RealmScoped:  true,
				}, []*ServiceOperation{
					{Name: "Find"},
					{Name: "List"},
					{
						Name: "ListAssociation",
						Params: map[string]interface{}{
							"Field":     "Provider",
							"Type":      "auth.v1.Provider",
							"Resource":  "providers",
							"ParamName": "provider",
						},
					},
					{
						Name: "Associate",
						Params: map[string]interface{}{
							"Field":     "Provider",
							"Resource":  "providers",
							"ParamName": "provider",
						},
					},
					{
						Name: "ListAssociation",
						Params: map[string]interface{}{
							"Field":     "Role",
							"Type":      "auth.v1.Role",
							"Resource":  "roles",
							"ParamName": "role",
						},
					},
					{
						Name: "Associate",
						Params: map[string]interface{}{
							"Field":     "Role",
							"Resource":  "roles",
							"ParamName": "role",
						},
					},
					{Name: "CreateUpdateDelete"},
					{Name: "UpdateConfig"},
				}),
				MakeServiceInput(&ServiceInput{
					Name:         "Role",
					Type:         "auth.v1.Role",
					Resource:     "role",
					TenantScoped: true,
					RealmScoped:  true,
				}, []*ServiceOperation{
					{Name: "Find"},
					{Name: "List"},
					{
						Name: "ListAssociation",
						Params: map[string]interface{}{
							"Field":     "Permission",
							"Type":      "auth.v1.Permission",
							"Resource":  "permissions",
							"ParamName": "permission",
						},
					},
					{
						Name: "Associate",
						Params: map[string]interface{}{
							"Field":     "Permission",
							"Resource":  "permissions",
							"ParamName": "permission",
						},
					},
					{Name: "CreateUpdateDelete"},
				}),
				MakeServiceInput(&ServiceInput{
					Name:         "Permission",
					Type:         "auth.v1.Permission",
					Resource:     "permission",
					TenantScoped: true,
					RealmScoped:  true,
				}, []*ServiceOperation{
					{Name: "Find"},
					{Name: "List"},
					{Name: "CreateUpdateDelete"},
				}),
				MakeServiceInput(&ServiceInput{
					Name:         "Scope",
					Type:         "auth.v1.Scope",
					Resource:     "scope",
					TenantScoped: true,
					RealmScoped:  true,
				}, []*ServiceOperation{
					{Name: "Find"},
					{Name: "List"},
					{
						Name: "ListAssociation",
						Params: map[string]interface{}{
							"Field":     "Permission",
							"Type":      "auth.v1.Permission",
							"Resource":  "permissions",
							"ParamName": "permission",
						},
					},
					{
						Name: "Associate",
						Params: map[string]interface{}{
							"Field":     "Permission",
							"Resource":  "permissions",
							"ParamName": "permission",
						},
					},
					{Name: "CreateUpdateDelete"},
				}),
			},
		},
	}

	for _, s := range input {
		writeTemplate(s, protoTemplate, fmt.Sprintf("./%s_grpc.proto", s.ResourceName))
	}

}

func writeTemplate(i *MainInput, templateContent, path string) {

	t := template.New("main")

	funcs := map[string]interface{}{}

	gen.AddFuncs(funcs, t)

	t.Funcs(funcs)

	t, err := t.Parse(templateContent)
	if err != nil {
		panic(err)
	}

	for k1, v1 := range operationTemplates {
		for k2, v2 := range map[string]string{
			"Proto":     v1.ProtoTemplate,
			"ProtoType": v1.ProtoTypeTemplate,
		} {
			t, err = t.New(k1 + k2).Parse(v2)
			if err != nil {
				panic(err)
			}
		}
	}

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	err = t.ExecuteTemplate(f, "main", i)
	if err != nil {
		panic(err)
	}
}
