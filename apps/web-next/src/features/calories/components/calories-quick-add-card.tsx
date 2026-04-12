import { ActionChipRow } from '../../../components/ui/ActionChipRow'
import { Panel } from '../../../components/ui/Panel'
import type { CaloriesQuickAddItem } from '../display'

export function CaloriesQuickAddCard(props: { items: CaloriesQuickAddItem[], submittingId?: string, onSelect: (item: CaloriesQuickAddItem) => void }) {
  return (
    <Panel title="Quick add">
      <ActionChipRow
        items={props.items.map((item) => ({
          ...item,
          disabled: props.submittingId === item.id,
          onClick: () => props.onSelect(item),
        }))}
        emptyMessage="No recent meals yet."
      />
    </Panel>
  )
}
