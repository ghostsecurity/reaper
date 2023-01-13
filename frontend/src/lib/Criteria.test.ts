import { test, expect } from 'vitest'
import { Comparison, Criteria, JoinType, Rule, Ruleset, Target } from './Criteria'
import { HttpRequest } from './Http'

test.each([
  {
    input: '',
    expected: new Ruleset(
      [],
      [],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'invalid query',
    expected: new Ruleset(
      [
        new Rule(Target.Raw, Comparison.CONTAINS, 'invalid query'),
      ],
      [],
      JoinType.AND,
    ),
    expectException: true,
  },
  {
    input: 'scheme eq http',
    expected: new Ruleset(
      [
        new Rule(Target.Scheme, Comparison.EQUAL, 'http'),
      ],
      [],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'scheme == http',
    expected: new Ruleset(
      [
        new Rule(Target.Scheme, Comparison.EQUAL, 'http'),
      ],
      [],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'scheme==http',
    expected: new Ruleset(
      [
        new Rule(Target.Scheme, Comparison.EQUAL, 'http'),
      ],
      [],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'scheme=="http"',
    expected: new Ruleset(
      [
        new Rule(Target.Scheme, Comparison.EQUAL, 'http'),
      ],
      [],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'scheme==\'http\'',
    expected: new Ruleset(
      [
        new Rule(Target.Scheme, Comparison.EQUAL, 'http'),
      ],
      [],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'scheme == "http"',
    expected: new Ruleset(
      [
        new Rule(Target.Scheme, Comparison.EQUAL, 'http'),
      ],
      [],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'scheme == "http" AND (host == "example.com" OR host == "example.org")',
    expected: new Ruleset(
      [
        new Rule(Target.Scheme, Comparison.EQUAL, 'http'),
      ],
      [
        new Ruleset(
          [
            new Rule(Target.Host, Comparison.EQUAL, 'example.com'),
            new Rule(Target.Host, Comparison.EQUAL, 'example.org'),
          ],
          [],
          JoinType.OR,
        ),
      ],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'scheme == "http" AND (host == "example.com" OR host == "example.org") AND path contains /foo',
    expected: new Ruleset(
      [
        new Rule(Target.Scheme, Comparison.EQUAL, 'http'),
        new Rule(Target.Path, Comparison.CONTAINS, '/foo'),
      ],
      [
        new Ruleset(
          [
            new Rule(Target.Host, Comparison.EQUAL, 'example.com'),
            new Rule(Target.Host, Comparison.EQUAL, 'example.org'),
          ],
          [],
          JoinType.OR,
        ),
      ],
      JoinType.AND,
    ),
    expectException: false,
  },
  {
    input: 'host is api.ghostbank.net AND (path contains api OR path contains auth) AND scheme is https',
    expected: new Ruleset(
      [
        new Rule(Target.Host, Comparison.EQUAL, 'api.ghostbank.net'),
        new Rule(Target.Scheme, Comparison.EQUAL, 'https'),
      ],
      [
        new Ruleset(
          [
            new Rule(Target.Path, Comparison.CONTAINS, 'api'),
            new Rule(Target.Path, Comparison.CONTAINS, 'auth'),
          ],
          [],
          JoinType.OR,
        ),
      ],
      JoinType.AND,
    ),
    expectException: false,
  },
])('Parse: $input', ({ input, expected, expectException }) => {
  const criteria = new Criteria(input)
  if (expectException) {
    expect(criteria.ParseError).not.toBeNull()
  } else {
    expect(criteria.ParseError).toBeNull()
  }
  expect(criteria.Raw).toEqual(input)
  expect(criteria.root).toEqual(expected)
})

test.each([
  {
    query: 'scheme == "http" AND (host == "example.com" OR host == "example.org") AND path contains /foo',
    url: 'http://example.com/foo',
    expected: true,
  },
  {
    query: 'scheme == "http" AND (host == "example.com" OR host == "example.org") AND path matches foo',
    url: 'http://example.com/foo',
    expected: true,
  },
  {
    query: 'host is api.ghostbank.net AND (path contains api OR path contains auth) AND scheme is https',
    url: 'https://api.ghostbank.net/api/v3/transactions',
    expected: true,
  },
  {
    query: 'host is api.ghostbank.net AND (path contains api OR path contains auth) AND scheme is https',
    url: 'https://api.ghostbonk.net/api/v3/transactions',
    expected: false,
  },
  {
    query: 'ghostbank',
    url: 'https://api.ghostbank.net/api/v3/transactions',
    expected: true,
    expectException: true,
  },
])('Match: $query -> $url', ({ query, url, expected, expectException }) => {
  const criteria = new Criteria(query)
  if (!expectException) {
    expect(criteria.ParseError).toBeNull()
  }
  const p = new URL(url)
  const protocol = p.protocol.replace(':', '')
  expect(criteria.Match(<HttpRequest>{
    Host: p.host,
    Path: p.pathname,
    Scheme: protocol,
    QueryString: p.search,
    Raw: `GET ${p.pathname} HTTP/1.1\r\nHost: ${p.host}\r\n\r\n`,
  })).toEqual(expected)
})
