import { For, createSignal } from "solid-js"

export type Pointer = {
	key: string | null
	at: number
}

export default function Ingredients() {
	let root: HTMLElement

	const [pointer, setPointer] = createSignal<Pointer>({ key: null, at: 0 })

	const [custom, setCustom] = createSignal([createSignal('')])
	const [generated, setGenerated] = createSignal([] as { key: string, ingredient: string[] }[])

	const setFocus = (pointer: Pointer) => {
		setPointer(pointer)
		getInput(pointer)?.focus()
	}

	const getInput = ({ key, at }: Pointer) => root.querySelector(`#${createId(key, at)}`) as HTMLInputElement | null
	const getCurrentInput = () => getInput(pointer())

	const previousPointer = () => {
		const { key, at } = pointer()
		if (at === 0) {
			// TODO
			return null
		}

		return { key, at: at - 1 }
	}

	const createNew = () => {
		const last = custom().at(-1);
		if (last && !last[0]()) {
			setTimeout(() => setFocus({ key: null, at: custom().length - 1 }), 0)
			return
		}

		const next = custom().length

		setCustom(customs => [...customs, createSignal('')])

		setTimeout(() => {
			setFocus({ key: null, at: next })
		}, 0)
	}

	const handleNavigation = (e: KeyboardEvent) => {
		if (e.key === "Enter") {
			createNew()
			return
		}

		if (e.key === 'Backspace') {
			if (!getCurrentInput()?.value) {
				const previous = previousPointer()

				setCustom(customs => {
					customs.splice(pointer().at, 1)
					return customs.slice()
				})

				if (previous)
					setTimeout(() => {
						setFocus(previous)
					}, 0)
			}

		}
	}

	const onCustomChange = (content: string, i: number) => {
		const [_, setCustom] = custom()[i]
		setCustom(content)
	}

	return <ul ref={el => { root = el }} onKeyDown={handleNavigation}>
		<li>
			<p>Custom</p>

			<ul>
				<For each={custom()}>
					{([ingredient], i) => <li>
						<input
							class="border-2 border-blue-500"
							id={createId(null, i())}
							type="text" value={ingredient()}
							onInput={e => onCustomChange(e.target.value, i())}
							onFocus={() => setPointer({ key: null, at: i() })}
						/>
					</li>}
				</For>
			</ul>
		</li>
	</ul>
}

function createId(key: string | null, index: number) {
	return `${key ?? 'custom'}-${index}`
}
