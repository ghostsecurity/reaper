import { HttpRequest } from './Http'

export {
  Criteria,
  Ruleset,
  Rule,
  Target,
  Comparison,
  JoinType,
}

class Criteria {
  Raw: string

  ParseError: Error | null = null

  root: Ruleset

  constructor(input: string) {
    this.Raw = input
    try {
      this.root = parse(input)
    } catch (e) {
      this.root = new Ruleset([
        new Rule(Target.Raw, Comparison.CONTAINS, input),
      ], [], JoinType.AND)
      this.ParseError = e as Error
    }
  }

  Match(request: HttpRequest): boolean {
    return this.root.Match(request)
  }
}

enum Comparison {
    EQUAL = 'eq',
    NOT_EQUAL = 'ne',
    CONTAINS = 'contains',
    MATCHES = 'matches',
}

const comparisonAliases = new Map<Comparison, string[]>([
  [Comparison.EQUAL, ['eq', '==', 'is']],
  [Comparison.NOT_EQUAL, ['neq', '!=']],
  [Comparison.CONTAINS, ['contains', 'includes', 'has', '*=']],
  [Comparison.MATCHES, ['matches', '~']],
])

enum Target {
    Scheme = 'scheme',
    Host = 'host',
    Path = 'path',
    Query = 'query',
    Raw = 'raw',
}

const targetAliases = new Map<Target, string[]>([
  [Target.Scheme, ['scheme', 'protocol', 'proto']],
  [Target.Host, ['hostname', 'host', 'domain']],
  [Target.Path, ['path']],
  [Target.Query, ['querystring', 'query', 'qs']],
])

enum JoinType {
    NONE = 'NONE',
    AND = 'AND',
    OR = 'OR',
}

const joinAliases = new Map<JoinType, string[]>([
  [JoinType.AND, ['and', '&&']],
  [JoinType.OR, ['or', '||']],
])

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
    case Target.Raw:
      field = req.Raw
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

class Ruleset {
  JoinType: JoinType = JoinType.AND

  Rulesets: Ruleset[]

  Rules: Rule[]

  constructor(rules: Rule[], rulesets: Ruleset[], join: JoinType) {
    this.Rules = rules
    this.Rulesets = rulesets
    this.JoinType = join
  }

  Match(req: HttpRequest): boolean {
    for (const rule of this.Rules) {
      const ruleResult = rule.Match(req)
      if (this.JoinType === JoinType.AND) {
        if (!ruleResult) {
          return false
        }
      } else if (ruleResult) {
        return true
      }
    }
    for (const ruleset of this.Rulesets) {
      const ruleResult = ruleset.Match(req)
      if (this.JoinType === JoinType.AND) {
        if (!ruleResult) {
          return false
        }
      } else if (ruleResult) {
        return true
      }
    }
    return this.JoinType === JoinType.AND
  }
}

class Reader {
  input: string

  pos: number

  saved: number

  constructor(input: string) {
    this.input = input
    this.pos = 0
    this.saved = 0
  }

  next(): string {
    const current = this.input[this.pos]
    this.pos++
    return current
  }

  peek(): string {
    return this.input[this.pos]
  }

  peekWord(): string {
    return this.input.slice(this.pos).trimStart().split(' ')[0]
  }

  readUntil(c: string): string {
    let str = ''
    while (this.peek() !== c && !this.complete()) {
      str += this.next()
    }
    return str
  }

  readWord(): string {
    this.skipWhitespace()
    let word = ''
    while (!this.complete() && this.peek() !== ' ' && this.peek() !== ')') {
      word += this.next()
    }
    return word
  }

  skipWhitespace(): boolean {
    let skipped = false
    while (this.peek() == ' ' && !this.complete()) {
      this.next()
      skipped = true
    }
    return skipped
  }

  complete(): boolean {
    return this.pos >= this.input.length
  }

  save(): void {
    this.saved = this.pos
  }

  restore(): void {
    this.pos = this.saved
  }

  prepend(c: string): void {
    this.input = c + this.input
  }
}

function parse(query: string): Ruleset {
  const reader = new Reader(query)
  return parseRuleset(reader, false)
}

function parseRuleset(reader: Reader, nested: boolean): Ruleset {
  const rules: Rule[] = []
  const rulesets: Ruleset[] = []
  let join: JoinType = JoinType.NONE
  let expectingRule = true
  while (!reader.complete()) {
    if (reader.skipWhitespace()) {
      continue
    }
    if (nested && reader.peek() === ')') {
      reader.next()
      expectingRule = false
      break
    }
    if (reader.peek() === '(') {
      reader.next()
      rulesets.push(parseRuleset(reader, true))
      expectingRule = false
      continue
    }
    if (expectingRule) {
      rules.push(parseRule(reader))
      expectingRule = false
      continue
    }

    let newJoin = JoinType.NONE
    joinAliases.forEach((aliases, type) => {
      aliases.forEach((alias) => {
        if (alias.toLowerCase() == reader.peekWord().toLowerCase()) {
          newJoin = <JoinType>type
          reader.readWord()
        }
      })
    })
    if (newJoin === JoinType.NONE) {
      throw new Error(`Expected either 'AND' or 'OR', found '${reader.peekWord()}'`)
    }
    if (join !== JoinType.NONE && join !== newJoin) {
      throw new Error(`Cannot mix ${join} and ${newJoin} without using brackets to group rules`)
    }
    join = newJoin
    expectingRule = true
  }
  if (join == JoinType.NONE) {
    join = JoinType.AND
  }
  return new Ruleset(rules, rulesets, join)
}

function parseRule(reader: Reader): Rule {
  reader.save()

  let target = Target.Raw
  let targetValid = false
  targetAliases.forEach((values, key) => {
    values.forEach((value) => {
      if (value.toLowerCase() === reader.peekWord().toLowerCase()) {
        target = <Target>key
        reader.readWord()
        targetValid = true
      }
    })
  })
  if (!targetValid) {
    targetAliases.forEach((values, key) => {
      values.forEach((value) => {
        if (reader.peekWord().toLowerCase().startsWith(value.toLowerCase())) {
          target = <Target>key
          const entireWord = reader.readWord()
          reader.prepend(entireWord.substring(value.length))
          targetValid = true
        }
      })
    })
  }
  if (!targetValid) {
    reader.restore()
    throw new Error(`invalid target '${reader.peekWord()}'`)
  }

  let comparison = Comparison.CONTAINS
  let comparisonValid = false
  comparisonAliases.forEach((values, key) => {
    values.forEach((value) => {
      if (value.toLowerCase() == reader.peekWord().toLowerCase()) {
        comparison = <Comparison>key
        reader.readWord()
        comparisonValid = true
      }
    })
  })
  if (!comparisonValid) {
    comparisonAliases.forEach((values, key) => {
      values.forEach((value) => {
        if (reader.peekWord().toLowerCase().startsWith(value.toLowerCase())) {
          comparison = <Comparison>key
          const entireWord = reader.readWord()
          reader.prepend(entireWord.substring(value.length))
          comparisonValid = true
        }
      })
    })
  }
  if (!comparisonValid) {
    reader.restore()
    throw new Error(`invalid comparison '${reader.peekWord()}'`)
  }

  reader.skipWhitespace()

  let value = ''
  switch (reader.peek()) {
  case '"':
    reader.next()
    value = reader.readUntil('"')
    reader.next()
    break
  case '\'':
    reader.next()
    value = reader.readUntil('\'')
    reader.next()
    break
  default:
    value = reader.readWord()
  }

  return new Rule(target, comparison, value)
}
