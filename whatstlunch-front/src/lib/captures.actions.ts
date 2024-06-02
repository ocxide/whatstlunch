import { createSignal, type Signal } from "solid-js";

export type Capture = {
	filename: string;
	ingredients: Signal<string[]>;
}

const [captures, setCaptures] = createSignal<Capture[]>([])

async function getIngredients(file: File) {
	const form = new FormData()
	form.append('image', file)

	const response = await fetch('http://192.168.0.8:3456/infer-ingredients', {
		method: 'POST',
		body: form
	})

	const data = await response.json()
	return data as string[]
}

const createCapture = (file: File) => ({
	filename: file.name,
	ingredients: createSignal<string[]>([])
})

export function insertCapture(file: File) {
	setCaptures(captures => {
		const i = captures.findIndex(c => c.filename === file.name)
		if (i === -1) {
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

export function removeCapture(filename: string) {
	setCaptures(captures => {
		const i = captures.findIndex(c => c.filename === filename)
		if (i === -1) return captures

		captures.splice(i, 1)
		return captures.slice()
	})
}

export { captures }
