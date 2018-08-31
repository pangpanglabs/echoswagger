package echoswagger

type (
	// Swagger represents an instance of a swagger object.
	// See https://swagger.io/specification/
	Swagger struct {
		Swagger             string                         `json:"swagger,omitempty"`
		Info                *Info                          `json:"info,omitempty"`
		Host                string                         `json:"host,omitempty"`
		BasePath            string                         `json:"basePath,omitempty"`
		Schemes             []string                       `json:"schemes,omitempty"`
		Consumes            []string                       `json:"consumes,omitempty"`
		Produces            []string                       `json:"produces,omitempty"`
		Paths               map[string]interface{}         `json:"paths"`
		Definitions         map[string]*JSONSchema         `json:"definitions,omitempty"`
		Parameters          map[string]*Parameter          `json:"parameters,omitempty"`
		Responses           map[string]*Response           `json:"responses,omitempty"`
		SecurityDefinitions map[string]*SecurityDefinition `json:"securityDefinitions,omitempty"`
		Tags                []*Tag                         `json:"tags,omitempty"`
		ExternalDocs        *ExternalDocs                  `json:"externalDocs,omitempty"`
	}

	// Info provides metadata about the API. The metadata can be used by the clients if needed,
	// and can be presented in the Swagger-UI for convenience.
	Info struct {
		Title          string                 `json:"title,omitempty"`
		Description    string                 `json:"description,omitempty"`
		TermsOfService string                 `json:"termsOfService,omitempty"`
		Contact        *Contact               `json:"contact,omitempty"`
		License        *License               `json:"license,omitempty"`
		Version        string                 `json:"version"`
		Extensions     map[string]interface{} `json:"-"`
	}

	// Contact contains the API contact information.
	Contact struct {
		// Name of the contact person/organization
		Name string `json:"name,omitempty"`
		// Email address of the contact person/organization
		Email string `json:"email,omitempty"`
		// URL pointing to the contact information
		URL string `json:"url,omitempty"`
	}

	// License contains the license information for the API.
	License struct {
		// Name of license used for the API
		Name string `json:"name,omitempty"`
		// URL to the license used for the API
		URL string `json:"url,omitempty"`
	}

	// Path holds the relative paths to the individual endpoints.
	Path struct {
		// Ref allows for an external definition of this path item.
		Ref string `json:"$ref,omitempty"`
		// Get defines a GET operation on this path.
		Get *Operation `json:"get,omitempty"`
		// Put defines a PUT operation on this path.
		Put *Operation `json:"put,omitempty"`
		// Post defines a POST operation on this path.
		Post *Operation `json:"post,omitempty"`
		// Delete defines a DELETE operation on this path.
		Delete *Operation `json:"delete,omitempty"`
		// Options defines a OPTIONS operation on this path.
		Options *Operation `json:"options,omitempty"`
		// Head defines a HEAD operation on this path.
		Head *Operation `json:"head,omitempty"`
		// Patch defines a PATCH operation on this path.
		Patch *Operation `json:"patch,omitempty"`
		// Parameters is the list of parameters that are applicable for all the operations
		// described under this path.
		Parameters []*Parameter `json:"parameters,omitempty"`
		// Extensions defines the swagger extensions.
		Extensions map[string]interface{} `json:"-"`
	}

	// Operation describes a single API operation on a path.
	Operation struct {
		// Tags is a list of tags for API documentation control. Tags can be used for
		// logical grouping of operations by resources or any other qualifier.
		Tags []string `json:"tags,omitempty"`
		// Summary is a short summary of what the operation does. For maximum readability
		// in the swagger-ui, this field should be less than 120 characters.
		Summary string `json:"summary,omitempty"`
		// Description is a verbose explanation of the operation behavior.
		// GFM syntax can be used for rich text representation.
		Description string `json:"description,omitempty"`
		// ExternalDocs points to additional external documentation for this operation.
		ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
		// OperationID is a unique string used to identify the operation.
		OperationID string `json:"operationId,omitempty"`
		// Consumes is a list of MIME types the operation can consume.
		Consumes []string `json:"consumes,omitempty"`
		// Produces is a list of MIME types the operation can produce.
		Produces []string `json:"produces,omitempty"`
		// Parameters is a list of parameters that are applicable for this operation.
		Parameters []*Parameter `json:"parameters,omitempty"`
		// Responses is the list of possible responses as they are returned from executing
		// this operation.
		Responses map[string]*Response `json:"responses,omitempty"`
		// Schemes is the transfer protocol for the operation.
		Schemes []string `json:"schemes,omitempty"`
		// Deprecated declares this operation to be deprecated.
		Deprecated bool `json:"deprecated,omitempty"`
		// Secury is a declaration of which security schemes are applied for this operation.
		Security []map[string][]string `json:"security,omitempty"`
		// Extensions defines the swagger extensions.
		Extensions map[string]interface{} `json:"-"`
	}

	// Parameter describes a single operation parameter.
	Parameter struct {
		// Name of the parameter. Parameter names are case sensitive.
		Name string `json:"name"`
		// In is the location of the parameter.
		// Possible values are "query", "header", "path", "formData" or "body".
		In string `json:"in"`
		// Description is`a brief description of the parameter.
		// GFM syntax can be used for rich text representation.
		Description string `json:"description,omitempty"`
		// Required determines whether this parameter is mandatory.
		Required bool `json:"required"`
		// Schema defining the type used for the body parameter, only if "in" is body
		Schema *JSONSchema `json:"schema,omitempty"`

		// properties below only apply if "in" is not body

		//  Type of the parameter. Since the parameter is not located at the request body,
		// it is limited to simple types (that is, not an object).
		Type string `json:"type,omitempty"`
		// Format is the extending format for the previously mentioned type.
		Format string `json:"format,omitempty"`
		// AllowEmptyValue sets the ability to pass empty-valued parameters.
		// This is valid only for either query or formData parameters and allows you to
		// send a parameter with a name only or an empty value. Default value is false.
		AllowEmptyValue bool `json:"allowEmptyValue,omitempty"`
		// Items describes the type of items in the array if type is "array".
		Items *Items `json:"items,omitempty"`
		// CollectionFormat determines the format of the array if type array is used.
		// Possible values are csv, ssv, tsv, pipes and multi.
		CollectionFormat string `json:"collectionFormat,omitempty"`
		// Default declares the value of the parameter that the server will use if none is
		// provided, for example a "count" to control the number of results per page might
		// default to 100 if not supplied by the client in the request.
		Default          interface{}   `json:"default,omitempty"`
		Maximum          *float64      `json:"maximum,omitempty"`
		ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
		Minimum          *float64      `json:"minimum,omitempty"`
		ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
		MaxLength        *int          `json:"maxLength,omitempty"`
		MinLength        *int          `json:"minLength,omitempty"`
		Pattern          string        `json:"pattern,omitempty"`
		MaxItems         *int          `json:"maxItems,omitempty"`
		MinItems         *int          `json:"minItems,omitempty"`
		UniqueItems      bool          `json:"uniqueItems,omitempty"`
		Enum             []interface{} `json:"enum,omitempty"`
		MultipleOf       float64       `json:"multipleOf,omitempty"`
		// Extensions defines the swagger extensions.
		Extensions map[string]interface{} `json:"-"`
	}

	// Response describes an operation response.
	Response struct {
		// Description of the response. GFM syntax can be used for rich text representation.
		Description string `json:"description,omitempty"`
		// Schema is a definition of the response structure. It can be a primitive,
		// an array or an object. If this field does not exist, it means no content is
		// returned as part of the response. As an extension to the Schema Object, its root
		// type value may also be "file".
		Schema *JSONSchema `json:"schema,omitempty"`
		// Headers is a list of headers that are sent with the response.
		Headers map[string]*Header `json:"headers,omitempty"`
		// Ref references a global API response.
		// This field is exclusive with the other fields of Response.
		Ref string `json:"$ref,omitempty"`
		// Extensions defines the swagger extensions.
		Extensions map[string]interface{} `json:"-"`
	}

	// Header represents a header parameter.
	Header struct {
		// Description is`a brief description of the parameter.
		// GFM syntax can be used for rich text representation.
		Description string `json:"description,omitempty"`
		//  Type of the header. it is limited to simple types (that is, not an object).
		Type string `json:"type,omitempty"`
		// Format is the extending format for the previously mentioned type.
		Format string `json:"format,omitempty"`
		// Items describes the type of items in the array if type is "array".
		Items *Items `json:"items,omitempty"`
		// CollectionFormat determines the format of the array if type array is used.
		// Possible values are csv, ssv, tsv, pipes and multi.
		CollectionFormat string `json:"collectionFormat,omitempty"`
		// Default declares the value of the parameter that the server will use if none is
		// provided, for example a "count" to control the number of results per page might
		// default to 100 if not supplied by the client in the request.
		Default          interface{}   `json:"default,omitempty"`
		Maximum          *float64      `json:"maximum,omitempty"`
		ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
		Minimum          *float64      `json:"minimum,omitempty"`
		ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
		MaxLength        *int          `json:"maxLength,omitempty"`
		MinLength        *int          `json:"minLength,omitempty"`
		Pattern          string        `json:"pattern,omitempty"`
		MaxItems         *int          `json:"maxItems,omitempty"`
		MinItems         *int          `json:"minItems,omitempty"`
		UniqueItems      bool          `json:"uniqueItems,omitempty"`
		Enum             []interface{} `json:"enum,omitempty"`
		MultipleOf       float64       `json:"multipleOf,omitempty"`
	}

	// SecurityDefinition allows the definition of a security scheme that can be used by the
	// operations. Supported schemes are basic authentication, an API key (either as a header or
	// as a query parameter) and OAuth2's common flows (implicit, password, application and
	// access code).
	SecurityDefinition struct {
		// Type of the security scheme. Valid values are "basic", "apiKey" or "oauth2".
		Type string `json:"type"`
		// Description for security scheme
		Description string `json:"description,omitempty"`
		// Name of the header or query parameter to be used when type is "apiKey".
		Name string `json:"name,omitempty"`
		// In is the location of the API key when type is "apiKey".
		// Valid values are "query" or "header".
		In string `json:"in,omitempty"`
		// Flow is the flow used by the OAuth2 security scheme when type is "oauth2"
		// Valid values are "implicit", "password", "application" or "accessCode".
		Flow string `json:"flow,omitempty"`
		// The oauth2 authorization URL to be used for this flow.
		AuthorizationURL string `json:"authorizationUrl,omitempty"`
		// TokenURL  is the token URL to be used for this flow.
		TokenURL string `json:"tokenUrl,omitempty"`
		// Scopes list the  available scopes for the OAuth2 security scheme.
		Scopes map[string]string `json:"scopes,omitempty"`
		// Extensions defines the swagger extensions.
		Extensions map[string]interface{} `json:"-"`
	}

	// Scope corresponds to an available scope for an OAuth2 security scheme.
	Scope struct {
		// Description for scope
		Description string `json:"description,omitempty"`
	}

	// ExternalDocs allows referencing an external resource for extended documentation.
	ExternalDocs struct {
		// Description is a short description of the target documentation.
		// GFM syntax can be used for rich text representation.
		Description string `json:"description,omitempty"`
		// URL for the target documentation.
		URL string `json:"url"`
	}

	// Items is a limited subset of JSON-Schema's items object. It is used by parameter
	// definitions that are not located in "body".
	Items struct {
		//  Type of the items. it is limited to simple types (that is, not an object).
		Type string `json:"type,omitempty"`
		// Format is the extending format for the previously mentioned type.
		Format string `json:"format,omitempty"`
		// Items describes the type of items in the array if type is "array".
		Items *Items `json:"items,omitempty"`
		// CollectionFormat determines the format of the array if type array is used.
		// Possible values are csv, ssv, tsv, pipes and multi.
		CollectionFormat string `json:"collectionFormat,omitempty"`
		// Default declares the value of the parameter that the server will use if none is
		// provided, for example a "count" to control the number of results per page might
		// default to 100 if not supplied by the client in the request.
		Default          interface{}   `json:"default,omitempty"`
		Maximum          *float64      `json:"maximum,omitempty"`
		ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
		Minimum          *float64      `json:"minimum,omitempty"`
		ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
		MaxLength        *int          `json:"maxLength,omitempty"`
		MinLength        *int          `json:"minLength,omitempty"`
		Pattern          string        `json:"pattern,omitempty"`
		MaxItems         *int          `json:"maxItems,omitempty"`
		MinItems         *int          `json:"minItems,omitempty"`
		UniqueItems      bool          `json:"uniqueItems,omitempty"`
		Enum             []interface{} `json:"enum,omitempty"`
		MultipleOf       float64       `json:"multipleOf,omitempty"`
	}

	// Tag allows adding meta data to a single tag that is used by the Operation Object. It is
	// not mandatory to have a Tag Object per tag used there.
	Tag struct {
		// Name of the tag.
		Name string `json:"name,omitempty"`
		// Description is a short description of the tag.
		// GFM syntax can be used for rich text representation.
		Description string `json:"description,omitempty"`
		// ExternalDocs is additional external documentation for this tag.
		ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
		// Extensions defines the swagger extensions.
		Extensions map[string]interface{} `json:"-"`
	}

	JSONSchema struct {
		Schema string `json:"$schema,omitempty"`
		// Core schema
		ID           string                 `json:"id,omitempty"`
		Title        string                 `json:"title,omitempty"`
		Type         JSONType               `json:"type,omitempty"`
		Items        *JSONSchema            `json:"items,omitempty"`
		Properties   map[string]*JSONSchema `json:"properties,omitempty"`
		Definitions  map[string]*JSONSchema `json:"definitions,omitempty"`
		Description  string                 `json:"description,omitempty"`
		DefaultValue interface{}            `json:"default,omitempty"`
		Example      interface{}            `json:"example,omitempty"`

		// Hyper schema
		Media    *JSONMedia `json:"media,omitempty"`
		ReadOnly bool       `json:"readOnly,omitempty"`
		Ref      string     `json:"$ref,omitempty"`
		XML      *XMLSchema `json:"xml,omitempty"`

		// Validation
		Enum                 []interface{} `json:"enum,omitempty"`
		Format               string        `json:"format,omitempty"`
		Pattern              string        `json:"pattern,omitempty"`
		Minimum              *float64      `json:"minimum,omitempty"`
		Maximum              *float64      `json:"maximum,omitempty"`
		MinLength            *int          `json:"minLength,omitempty"`
		MaxLength            *int          `json:"maxLength,omitempty"`
		Required             []string      `json:"required,omitempty"`
		AdditionalProperties *JSONSchema   `json:"additionalProperties,omitempty"`

		// Union
		AnyOf []*JSONSchema `json:"anyOf,omitempty"`
	}

	// JSONType is the JSON type enum.
	JSONType string

	// JSONMedia represents a "media" field in a JSON hyper schema.
	JSONMedia struct {
		BinaryEncoding string `json:"binaryEncoding,omitempty"`
		Type           string `json:"type,omitempty"`
	}

	// JSONLink represents a "link" field in a JSON hyper schema.
	JSONLink struct {
		Title        string      `json:"title,omitempty"`
		Description  string      `json:"description,omitempty"`
		Rel          string      `json:"rel,omitempty"`
		Href         string      `json:"href,omitempty"`
		Method       string      `json:"method,omitempty"`
		Schema       *JSONSchema `json:"schema,omitempty"`
		TargetSchema *JSONSchema `json:"targetSchema,omitempty"`
		MediaType    string      `json:"mediaType,omitempty"`
		EncType      string      `json:"encType,omitempty"`
	}

	XMLSchema struct {
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
		Prefix    string `json:"prefix,omitempty"`
		Attribute string `json:"attribute,omitempty"`
		Wrapped   bool   `json:"wrapped,omitempty"`
	}
)
