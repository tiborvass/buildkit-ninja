package ninja2llb

import (
	"strings"
)

func Expand(s string, getter Getter) string {
	var sb strings.Builder
	foundDollar := false
	inVariableName := false
	openBrace := false
	j := 0
	for i := 0; i < len(s); i++ {
		switch ch := s[i]; ch {
		case '$':
			if openBrace {
				panic("invalid '$' inside variable name")
			}
			if inVariableName {
				v, _ := getter.Get(s[j:i])
				sb.WriteString(v)
				inVariableName = false
				j = i
			}
			foundDollar = true
			sb.WriteString(s[j:i])
			i++ // check character after '$'
			if i == len(s) {
				panic("invalid trailing '$'")
			}
			j = i
			switch s[i] {
			case '$', ':', ' ', '\n':
				sb.WriteByte(s[i])
				j++
			case '{':
				if i+1 == len(s) {
					panic("invalid trailing '${'")
				}
				j++
				openBrace = true
			default:
				inVariableName = true
			}
		case '}':
			if openBrace {
				v, _ := getter.Get(s[j:i])
				sb.WriteString(v)
				openBrace = false
				j = i + 1
			}
		default:
			if inVariableName {
				if (ch >= 'a' && ch <= 'z') ||
					(ch >= 'A' && ch <= 'Z') ||
					(ch >= '0' && ch <= '9') ||
					ch == '_' {
					continue
				}
				v, _ := getter.Get(s[j:i])
				sb.WriteString(v)
				inVariableName = false
				j = i
			}
		}
	}
	if !foundDollar {
		return s
	}
	if openBrace {
		panic("unclosed variable substitution '${'")
	}
	if inVariableName {
		v, _ := getter.Get(s[j:])
		sb.WriteString(v)
	} else {
		sb.WriteString(s[j:])
	}
	return sb.String()
}
