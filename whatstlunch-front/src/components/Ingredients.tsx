import { For, createEffect, createSignal } from "solid-js"

export default function Ingredients() {
	let root: HTMLElement

	const [pointer, setPointer] = createSignal({ key: null, at: 0 })

	const [custom, setCustom] = createSignal([''])
	const [generated, setGenerated] = createSignal([] as { key: string, ingredient: string[] }[])

	const getCurrentInput = () => {
		const { key, at } = pointer()
		return root.querySelector(`#${createId(key, at)}`) as HTMLInputElement | null;
	}

	const previousPointer = () => {
		const { key, at } = pointer()
		if (at === 0) {
			// TODO
			return null
		}

		return { key, at: at - 1 }
	}

	const handleNavigation = (e: KeyboardEvent) => {
		console.log(e.key)

		if (e.key === "Enter") {
			const next = custom().length

			setCustom(customs => [...customs, ''])

			setTimeout(() => {
				setPointer({ key: null, at: next })
			}, 100)
		}

		if (e.key === 'Backspace') {
			if (!getCurrentInput()?.value) {
				const previous = previousPointer()

				setCustom(customs => {

					return customs
				})
			}

		}
	}

	createEffect(() => {
		getCurrentInput()?.focus()
	})

	const onCustomChange = (content: string, i: number) => {
		setCustom(customs => {
			customs[i] = content
			return customs
		})
	}

	return <ul ref={el => { root = el }} onKeyDown={handleNavigation}>
		<li>
			<p>Custom</p>

			<ul>
				<For each={custom()}>
					{(ingredient, i) => <li>
						<input id={createId(null, i())} type="text" value={ingredient} onChange={e => onCustomChange(e.target.value, i())} />
					</li>}
				</For>
			</ul>
		</li>
	</ul>
}

function createId(key: string | null, index: number) {
	return `${key ?? 'custom'}-${index}`
}
