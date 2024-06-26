import { Status, captures } from "@/lib/captures.actions"
import { customStore } from "@/lib/custom-ingredients.store"
import { For, Show, createSignal } from "solid-js"

export type Pointer = {
	key: string | null
	at: number
}

export default function Ingredients() {
	const [pointer, setPointer] = createSignal<Pointer>({ key: null, at: 0 })
	const [customs, setCustom] = customStore

	const setFocus = (pointer: Pointer) => {
		setPointer(pointer)
		getInput(pointer)?.focus()
	}

	const getInput = ({ key, at }: Pointer) => document.getElementById(createId(key, at)) as HTMLInputElement | null
	const getCurrentInput = () => getInput(pointer())

	const createNew = () => {
		const last = customs().at(-1);
		if (last && !last[0]()) {
			setTimeout(() => setFocus({ key: null, at: customs().length - 1 }), 0)
			return
		}

		const next = customs().length

		setCustom(customs => [...customs, createSignal('')])

		setTimeout(() => {
			setFocus({ key: null, at: next })
		}, 0)
	}

	const handleCustomBackspace = () => {
		if (pointer().key != null || customs().length < 2 || getCurrentInput()?.value) return

		let previous: Pointer | null = null
		if (pointer().at > 0)
			previous = { key: null, at: pointer().at - 1 }

		setCustom(customs => {
			customs.splice(pointer().at, 1)
			return customs.slice()
		})

		if (previous)
			setTimeout(() => {
				setFocus(previous)
			}, 0)
	}

	const handleNavigation = (e: KeyboardEvent) => {
		if (e.key === "Enter") {
			createNew()
			return
		}

		if (e.key === 'Backspace') {
			handleCustomBackspace()
		}
	}

	const onCustomChange = (content: string, i: number) => {
		const [_, setCustom] = customs()[i]
		setCustom(content)
	}

	const onCapturedChange = (content: string, key: string, i: number) => {
		const capture = captures().find(c => c.filename === key)
		if (!capture) return
		const [_, set] = capture.ingredients

		set(ingredients => {
			ingredients.splice(i, 1)
			return ingredients.slice()
		})

		const last = customs().at(-1);
		if (last && !last[0]()) {
			const [_, setLast] = last
			setLast(content)
		}
		else {
			setCustom(customs => [...customs, createSignal(content)])
		}

		setTimeout(() => setFocus({ key: null, at: customs().length - 1 }), 0)
	}

	return <ul class="flex flex-col gap-2" onKeyDown={handleNavigation}>
		<For each={captures()}>
			{(capture) => (<li>
				<p class="font-bold">{capture.filename}</p>

				<Show when={capture.status[0]() === Status.Error}>
					<p class="text-red-500">Error</p>
				</Show>

				<Show when={capture.status[0]() === Status.Loading}>
					<p class="text-gray-300 font-light">Loading...</p>
				</Show>

				<ul>
					<For each={capture.ingredients[0]()}>
						{(ingredient, i) => <li>
							<input
								class="border-2 border-blue-500"
								id={createId(capture.filename, i())}
								type="text" value={ingredient}
								onInput={e => onCapturedChange(e.target.value, capture.filename, i())}
								onFocus={() => setPointer({ key: capture.filename, at: i() })}
							/>
						</li>}
					</For>
				</ul>
			</li>)}
		</For>

		<li>
			<p>Custom</p>

			<ul>
				<For each={customs()}>
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
