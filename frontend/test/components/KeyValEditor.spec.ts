import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import KeyValEditor from '../../src/components/KeyValEditor.vue'
import AutocompleteInput from '../../src/components/Shared/AutocompleteInput.vue'

describe('KeyValEditor', () => {
  test('Shows empty message and no table/rows when there are no key values and read-only', () => {
    const wrapper = mount(KeyValEditor, {
      propsData: {
        data: [],
        readonly: true,
        emptyMessage: 'No headers found for this request.',
      },
    })

    expect(wrapper.text()).contains('No headers found for this request.')
    expect(wrapper.findAll('table').length).toBe(0)
    expect(wrapper.findAll('tr').length).toBe(0)
  })
  test('Shows single empty row when there are no key values and not read-only', () => {
    const wrapper = mount(KeyValEditor, {
      propsData: {
        data: [],
        readonly: false,
        emptyMessage: 'should not show',
      },
    })

    expect(wrapper.text()).not.contains('should not show')
    expect(wrapper.findAll('table').length).toBe(1)
    expect(wrapper.findAll('tr').length).toBe(1)
    expect((wrapper.get('td:nth-child(1) input').element as HTMLInputElement).readOnly).toBe(false)
    expect((wrapper.get('td:nth-child(2) input').element as HTMLInputElement).readOnly).toBe(false)
  })
  test('Delete button should show on all rows except final', () => {
    const wrapper = mount(KeyValEditor, {
      propsData: {
        data: [
          {
            Key: 'a',
            Value: 'b',
          },
        ],
        readonly: false,
      },
    })
    expect(wrapper.findAll('table tr:nth-child(1) td:nth-child(3) a').length).toBe(1)
    expect(wrapper.findAll('table tr:nth-child(2) td:nth-child(3) a').length).toBe(0)
  })
  test('Shows single row (read-only)', () => {
    const wrapper = mount(KeyValEditor, {
      propsData: {
        data: [
          {
            Key: 'a',
            Value: 'b',
          },
        ],
        readonly: true,
        emptyMessage: 'should not show',
      },
    })

    expect(wrapper.text()).not.contains('should not show')
    expect(wrapper.findAll('table tr').length).toBe(1)
    expect(wrapper.findAll('table tr td').length).toBe(2) // no delete buttons

    const propsKey = wrapper.get('td:nth-child(1)').findComponent(AutocompleteInput).props()
    const propsValue = wrapper.get('td:nth-child(2)').findComponent(AutocompleteInput).props()

    expect(propsKey.value).toBe('a')
    expect(propsKey.readonly).toBe(true)

    expect(propsValue.value).toBe('b')
    expect(propsValue.readonly).toBe(true)
  })
  test('Shows single row with extra blank row (not read-only)', () => {
    const wrapper = mount(KeyValEditor, {
      propsData: {
        data: [
          {
            Key: 'a',
            Value: 'b',
          },
        ],
        readonly: false,
        emptyMessage: 'should not show',
      },
    })

    expect(wrapper.text()).not.contains('should not show')
    expect(wrapper.findAll('table tr').length).toBe(2)
    expect(wrapper.findAll('table tr:nth-child(1) td').length).toBe(3) // with delete button
    expect((wrapper.get('table tr:nth-child(1) td:nth-child(1) input').element as HTMLInputElement).value).toBe('a')
    expect((wrapper.get('table tr:nth-child(1) td:nth-child(2) input').element as HTMLInputElement).value).toBe('b')
    expect((wrapper.get('table tr:nth-child(2) td:nth-child(1) input').element as HTMLInputElement).value).toBe('')
    expect((wrapper.get('table tr:nth-child(2) td:nth-child(2) input').element as HTMLInputElement).value).toBe('')
  })
  test('Delete button deletes row', () => {
    const wrapper = mount(KeyValEditor, {
      propsData: {
        data: [
          {
            Key: 'a',
            Value: 'first',
          },
          {
            Key: 'b',
            Value: 'second',
          },
          {
            Key: 'c',
            Value: 'last',
          },
        ],
        readonly: false,
      },
    })

    // delete 'b'
    wrapper.get('tr:nth-child(2) td:nth-child(3) a').trigger('click').then(() => {
      expect(wrapper.emitted()).toHaveProperty('publish')
      expect(wrapper.emitted().publish).toHaveLength(1)
      expect(wrapper.emitted().publish[0]).toStrictEqual([
        [
          {
            Key: 'a',
            Value: 'first',
          },
          {
            Key: 'c',
            Value: 'last',
          },
        ],
      ])
    })
  })
  test('Editing row causes it to be published', () => {
    const wrapper = mount(KeyValEditor, {
      propsData: {
        data: [
          {
            Key: 'a',
            Value: 'first',
          },
          {
            Key: 'b',
            Value: 'second',
          },
          {
            Key: 'c',
            Value: 'last',
          },
        ],
        readonly: false,
      },
    })

    // edit 'b' from 'second' to 'changed'
    wrapper.get('tr:nth-child(2) td:nth-child(2) input').setValue('changed').then(
      () => {
        expect(wrapper.emitted()).toHaveProperty('publish')
        expect(wrapper.emitted().publish).toHaveLength(1)
        expect(wrapper.emitted().publish[0]).toStrictEqual([
          [
            {
              Key: 'a',
              Value: 'first',
            },
            {
              Key: 'b',
              Value: 'changed',
            },
            {
              Key: 'c',
              Value: 'last',
            },
          ],
        ])
      },
    ).catch(
      err => {
        expect(err).toBe(null)
      },
    )
  })
})
