import { CheckToggleButton } from '../../../components/ui/ActionButton'
import { List, ListItem } from '../../../components/ui/List'
import { Panel } from '../../../components/ui/Panel'
import { CounterPill } from '../../../components/ui/Pill'
import type { ExpenseChecklistItem } from '../../../core/api/types'
import { formatEuro } from '../../../core/format/money'

export function ObligationsCard(props: { obligations: ExpenseChecklistItem[], busyId?: string, onComplete: (item: ExpenseChecklistItem) => void }) {
  return (
    <Panel
      title="Obligations"
      description={props.obligations.length > 0 ? <CounterPill value={props.obligations.length} label="Left" /> : undefined}
      collapsibleId="expenses_obligations"
    >
      <List emptyMessage="All obligations paid for this month!" variant="flush">
        {props.obligations.map((item) => (
          <ListItem
            id={item.id}
            title={item.title}
            value={`-${formatEuro(item.amount)}`}
            valueTone="danger"
            valueStyle="mono"
            prefix={<CheckToggleButton disabled={props.busyId === item.id} onClick={() => props.onComplete(item)} />}
          />
        ))}
      </List>
    </Panel>
  )
}
