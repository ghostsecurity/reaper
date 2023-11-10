import { JoinType, Rule } from './Rule'
import { HttpRequest } from '../api/packaging'

export default class Ruleset {
  JoinType: JoinType = JoinType.AND

    Rulesets: Ruleset[] /* eslint-disable-line */

  Rules: Rule[]

  constructor(rules: Rule[], rulesets: Ruleset[], join: JoinType) {
    this.Rules = rules
    this.Rulesets = rulesets
    this.JoinType = join
  }

  Match(req: HttpRequest): boolean {
    const NoMatch = {}
    const Match = {}
    try {
      this.Rules.forEach(rule => {
        const ruleResult = rule.Match(req)
        if (this.JoinType === JoinType.AND) {
          if (!ruleResult) {
            throw NoMatch
          }
        } else if (ruleResult) {
          throw Match
        }
      })
      this.Rulesets.forEach(ruleset => {
        const ruleResult = ruleset.Match(req)
        if (this.JoinType === JoinType.AND) {
          if (!ruleResult) {
            throw NoMatch
          }
        } else if (ruleResult) {
          throw Match
        }
      })
    } catch (e) {
      if (e === Match) return true
      if (e === NoMatch) return false
      throw e
    }
    return this.JoinType === JoinType.AND
  }
}
