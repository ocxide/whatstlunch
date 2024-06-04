import { createSignal, type Signal } from "solid-js";

export const enum Status {
	Loading,
	Error,
	Ok
}

export type Capture = {
	filename: string;
	ingredients: Signal<string[]>;
	status: Signal<Status>;
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

const createCapture = (file: File): Capture => ({
	filename: file.name,
	ingredients: createSignal<string[]>([]),
	status: createSignal<Status>(Status.Loading),
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
		const capture = captures().find(c => c.filename === file.name);
		if (!capture) return

		const { status, ingredients: ingredientsSignal } = capture

		{
			const [_, setStatus] = status
			setStatus(Status.Ok)
		}
		{
			const [_, setIngredients] = ingredientsSignal
			setIngredients(ingredients)
		}

	}).catch(() => {
		const status = captures().find(c => c.filename === file.name)?.status;
		if (!status) return

		const [_, set] = status
		set(Status.Error)
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
