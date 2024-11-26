import moment from 'moment'

const userLocale = navigator.language || 'en-US'

interface UtilStore {
  customNumberFormat: (number: number) => string
  customNumberFormatK: (number: number) => string
  localeTimestamp: (ts: number) => string
  localeTimestampShort: (ts: number) => string
  localeDate: (ts: number) => string
  timeAgoInWords: (date: Date) => string
}

export function useUtilStore(): UtilStore {
  const localeTimestamp = (ts: number) => {
    const date = new Date(ts * 1000)
    const timeString = date.toLocaleTimeString(userLocale, {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    })
    return timeString
  }

  const localeTimestampShort = (ts: number) => {
    const date = new Date(ts * 1000)
    const timeString = date.toLocaleTimeString(userLocale, {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    })
    return timeString
  }

  const localeDate = (ts: number) => {
    const date = new Date(ts * 1000)
    const dateString = date.toLocaleDateString(userLocale, {
      month: 'short',
      day: 'numeric',
    })
    return dateString
  }

  const timeAgoInWords = (date: Date) => {
    return moment(date).fromNow()
  }

  const customNumberFormat = (n: number) => {
    const formatter = new Intl.NumberFormat(userLocale, {
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    })

    return formatter.format(n);
  }

  const customNumberFormatK = (n: number) => {
    const formatter = new Intl.NumberFormat(userLocale, {
      minimumFractionDigits: 0,
      maximumFractionDigits: 1
    })

    if (n >= 1000) {
      // Format the number as "K" (thousands) when it's 1000 or greater
      return formatter.format(n / 1000) + 'K';
    }
    return formatter.format(n)
  }

  return {
    customNumberFormat,
    customNumberFormatK,
    localeTimestamp,
    localeTimestampShort,
    localeDate,
    timeAgoInWords
  }
}
