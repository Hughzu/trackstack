import { Panel } from './Panel'
import { Pill } from './Pill'

type FeaturePlaceholderProps = {
  title: string
  route: string
  description: string
  bullets: string[]
}

export function FeaturePlaceholder(props: FeaturePlaceholderProps) {
  return (
    <Panel title={props.title} description={props.description} href={props.route}>
      <div class="flex flex-1 flex-col gap-4">
        <Pill tone="accent">Scaffolded</Pill>
        <div class="rounded-xl border border-border bg-panel px-3 py-2 text-xs text-text-muted">Route: {props.route}</div>
        <ul class="text-sm leading-6 text-text-main">
          {props.bullets.map((bullet) => (
            <li>{bullet}</li>
          ))}
        </ul>
      </div>
    </Panel>
  )
}
