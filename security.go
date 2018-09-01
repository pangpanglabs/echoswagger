package echoswagger

import "errors"

type SecurityType string

const (
	SecurityBasic  SecurityType = "basic"
	SecurityOAuth2 SecurityType = "oauth2"
	SecurityAPIKey SecurityType = "apiKey"
)

type SecurityInType string

const (
	SecurityInQuery  SecurityInType = "query"
	SecurityInHeader SecurityInType = "header"
)

type OAuth2FlowType string

const (
	OAuth2FlowImplicit    OAuth2FlowType = "implicit"
	OAuth2FlowPassword    OAuth2FlowType = "password"
	OAuth2FlowApplication OAuth2FlowType = "application"
	OAuth2FlowAccessCode  OAuth2FlowType = "accessCode"
)

func (r *Root) checkSecurity(name string) bool {
	if name == "" {
		return false
	}
	if _, ok := r.spec.SecurityDefinitions[name]; ok {
		return false
	}
	return true
}

func setSecurity(security []map[string][]string, names ...string) []map[string][]string {
	m := make(map[string][]string)
	for _, name := range names {
		m[name] = make([]string, 0)
	}
	return append(security, m)
}

func setSecurityWithScope(security []map[string][]string, s ...map[string][]string) []map[string][]string {
	for _, t := range s {
		if len(t) == 0 {
			continue
		}
		for k, v := range t {
			if len(v) == 0 {
				t[k] = make([]string, 0)
			}
		}
		security = append(security, t)
	}
	return security
}

func (o *Operation) addSecurity(defs map[string]*SecurityDefinition, security []map[string][]string) error {
	for _, scy := range security {
		for k := range scy {
			if _, ok := defs[k]; !ok {
				return errors.New("echoswagger: not found SecurityDefinition with name: " + k)
			}
		}
		if containsMap(o.Security, scy) {
			continue
		}
		o.Security = append(o.Security, scy)
	}
	return nil
}
