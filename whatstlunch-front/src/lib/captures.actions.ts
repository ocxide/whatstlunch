import { createSignal, type Signal } from "solid-js";

type Capture = {
	filename: string;
	ingredients: Signal<string[]>;
}

const [captures, setCaptures] = createSignal<Capture[]>([])

const getIngredients = (_file: File) => new Promise<string[]>(resolve => {
	setTimeout(() => {
		resolve(['apple', 'banana'])
	}, 2000)
})

const createCapture = (file: File) => ({
	filename: file.name,
	ingredients: createSignal<string[]>([])
})

export function insertCapture(file: File) {
	setCaptures(captures => {
		const i = captures.findIndex(c => c.filename === file.name)
		if (i !== -1) {
			return [...captures, createCapture(file)]
		}

		captures[i] = createCapture(file)
		return captures.slice()
	})

	getIngredients(file).then(ingredients => {
		const signal = captures().find(c => c.filename === file.name)?.ingredients;
		if (!signal) return

		const [_, set] = signal
		set(ingredients)
	})
}

export { captures }
