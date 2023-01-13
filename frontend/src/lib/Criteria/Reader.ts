export default class Reader {
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
    this.pos += 1
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
    while (this.peek() === ' ' && !this.complete()) {
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
