import { HttpRequest } from '../Http'

enum Comparison {
  EQUAL = 'eq',
  NOT_EQUAL = 'ne',
  CONTAINS = 'contains',
  MATCHES = 'matches',
}

enum Target {
  Scheme = 'scheme',
  Host = 'host',
  Path = 'path',
  Query = 'query',
  Body = 'body',
}

enum JoinType {
  NONE = 'NONE',
  AND = 'AND',
  OR = 'OR',
}

class Rule {
  Target: Target

  Comparison: Comparison

  Value: string

  constructor(target: Target, comparison: Comparison, value: string) {
    this.Target = target
    this.Comparison = comparison
    this.Value = value
  }

  Match(req: HttpRequest): boolean {
    let field = ''
    switch (this.Target) {
    case Target.Scheme:
      field = req.Scheme
      break
    case Target.Host:
      field = req.Host
      break
    case Target.Path:
      field = req.Path
      break
    case Target.Query:
      field = req.QueryString
      break
    case Target.Body:
      field = req.Body
      break
    default:
      return false
    }
    switch (this.Comparison) {
    case Comparison.EQUAL:
      return field === this.Value
    case Comparison.NOT_EQUAL:
      return field !== this.Value
    case Comparison.CONTAINS:
      return field.indexOf(this.Value) > -1
    case Comparison.MATCHES:
      try {
        return field.match(this.Value) !== null
      } catch (e) {
        return false
      }
    default:
      return false
    }
  }
}

export { Comparison, Target, JoinType, Rule }
