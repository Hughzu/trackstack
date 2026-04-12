const countFormatter = new Intl.NumberFormat('en-IE', {
  maximumFractionDigits: 1,
})

export function formatCount(value: number) {
  return countFormatter.format(value)
}
