import { Target, Comparison, JoinType, Rule } from './Rule'
import Ruleset from './Ruleset'
import Reader from './Reader'
import { HttpRequest } from '../Http'

export { Criteria, Ruleset }

class Criteria {
  Raw: string

  ParseError: Error | null = null

  root: Ruleset

  constructor(input: string) {
    this.Raw = input
    try {
      this.root = parse(input)
    } catch (e) {
      this.root = new Ruleset([new Rule(Target.Raw, Comparison.CONTAINS, input)], [], JoinType.AND)
      this.ParseError = e as Error
    }
  }

  Match(request: HttpRequest): boolean {
    return this.root.Match(request)
  }
}

const comparisonAliases = new Map<Comparison, string[]>([
  [Comparison.EQUAL, ['eq', '==', 'is']],
  [Comparison.NOT_EQUAL, ['neq', '!=']],
  [Comparison.CONTAINS, ['contains', 'includes', 'has', '*=']],
  [Comparison.MATCHES, ['matches', '~']],
])

const targetAliases = new Map<Target, string[]>([
  [Target.Scheme, ['scheme', 'protocol', 'proto']],
  [Target.Host, ['hostname', 'host', 'domain']],
  [Target.Path, ['path']],
  [Target.Query, ['querystring', 'query', 'qs']],
])

const joinAliases = new Map<JoinType, string[]>([
  [JoinType.AND, ['and', '&&']],
  [JoinType.OR, ['or', '||']],
])

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
        if (alias.toLowerCase() === reader.peekWord().toLowerCase()) {
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
  if (join === JoinType.NONE) {
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
      if (value.toLowerCase() === reader.peekWord().toLowerCase()) {
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
